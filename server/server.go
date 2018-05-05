package server

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/wavyllama/chat/core"
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/protocol"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"time"
	"net/http"
	"os/exec"
)

const (
	Port    uint16 = 4242
	localhost = "localhost"
)

// Server holds the user and all of his sessions
type Server struct {
	User       *db.User
	LastFriend *db.Friend
	Listener   *http.Server
	Tunnel *exec.Cmd
	Sessions   *[]*Session

	onReceiveFriendMessage func(m *FriendMessage)
	onAcceptFriend         func(displayName string)
	onReceiveChatMessage   func(message []byte, friend *db.Friend, time time.Time)
	onInfoReceive          func(messageToDisplay string)
}

func init() {
	gob.Register(&FriendMessage{})
	gob.Register(&HandshakeMessage{})
	gob.Register(&ChatMessage{})
}

func InitServer(user *db.User) (*Server) {
	server := Server{}
	// Init no-op function handlers
	server.onReceiveFriendMessage = func(m *FriendMessage) {}
	server.onAcceptFriend = func(displayName string) {}
	server.onReceiveChatMessage = func(message []byte, friend *db.Friend, time time.Time) {}
	server.onInfoReceive = func(messageToDisplay string) {}

	// Initialize the session struct to a pointer
	var sessions []*Session
	server.Sessions = &sessions

	server.User = user
	server.LastFriend = new(db.Friend)

	// Updates the IP address of the user and create a friend for yourself
	server.User.DeleteFriend(db.Self)
	server.User.AddFriend(db.Self, user.MAC, localhost, user.Username)
	return &server
}

func (s *Server) InitUIHandlers(onReceiveFriendMessage func(m *FriendMessage),
								onAcceptFriend func(displayName string),
								onReceiveChatMessage func(message []byte, friend *db.Friend, time time.Time),
								onInfoReceive func(messageToDisplay string)) {
	s.onReceiveFriendMessage = onReceiveFriendMessage
	s.onAcceptFriend = onAcceptFriend
	s.onReceiveChatMessage = onReceiveChatMessage
	s.onInfoReceive = onInfoReceive
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message

	var wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := wsupgrader.Upgrade(w, r, nil)

	_, reader, err := conn.NextReader()
	if err != nil {
		panic(err)
	}

	dec := gob.NewDecoder(reader)

	if err = dec.Decode(&msg); err != nil {
		panic(err)
	}

	sourceMAC, _, sourceUsername := msg.SourceID()
	if sourceMAC == "" || sourceUsername == "" {
		log.Panicln("Received ill-formatted message")
	}

	if msg.DestID() != s.User.Username {
		log.Panicln("Received a message but it was not for me.")
		return
	}
	sessions := s.GetSessionsWithFriend(sourceMAC, sourceUsername)
	friend := s.User.GetFriendByUsernameAndMAC(sourceUsername, sourceMAC)

	switch msg.(type) {
	case *FriendMessage:
		if friend != nil {
			return
		}
		s.handleFriendMessage(msg.(*FriendMessage))
	case *HandshakeMessage:
		s.handleHandshakeMessage(friend, msg.(*HandshakeMessage))
	case *ChatMessage:
		if len(sessions) == 0 {
			return
		}
		s.handleChatMessage(msg.(*ChatMessage))
	}
}

// Handles a friend message
func (s *Server) handleFriendMessage(msg *FriendMessage) {
	// Set the last friend request you received
	s.LastFriend.MAC, s.LastFriend.IP, s.LastFriend.Username = msg.SourceID()
	// Display it on the UI
	s.onReceiveFriendMessage(msg)
}

func (s *Server) AcceptedFriend(displayName string) {
	lower := strings.ToLower(displayName)
	if lower == db.Self || s.User.IsFriendsWith(lower) {
		// TODO: error
	} else {
		s.User.AddFriend(displayName, s.LastFriend.MAC, s.LastFriend.IP, s.LastFriend.Username)
		s.SendFriendRequest(s.LastFriend.IP, s.LastFriend.Username)
		s.onAcceptFriend(displayName)
	}
}

// Handles a handshake message
func (s *Server) handleHandshakeMessage(friend *db.Friend, msg *HandshakeMessage) {
	sourceMAC, sourceIP, sourceUsername := msg.SourceID()
	messageYourself := sourceMAC == s.User.MAC && sourceUsername == s.User.Username
	sessions := s.GetSessionsWithFriend(sourceMAC, sourceUsername)

	// We are in a handshake, so the friend should exist already
	if friend == nil {
		return
	}
	var sess *Session
	round := msg.Round
	protoType := msg.ProtoType

	// In a handshake, create a new session if there aren't the required number of sessions in either situation
	if len(sessions) != 2 && messageYourself || (len(sessions) != 1 && !messageYourself) {
		sess = NewSessionFromUserAndMessage(s.User, friend, protoType)
		*(*s).Sessions = append(*(*s).Sessions, sess)
	} else if len(sessions) == 2 && messageYourself {
		// Communicating between yourself, rotate sessions based on round (even/odd)
		sess = sessions[round%2]
	} else {
		sess = sessions[0]
	}

	dec, err := sess.Proto.Decrypt(msg.Secret, s.onInfoReceive)
	switch errorType := err.(type) {
	case protocol.OTRHandshakeStep:
		// Send each part of the handshake message back and immediately return
		for _, stepMessage := range dec {
			reply := new(HandshakeMessage)
			sourceIP = s.User.IP
			if messageYourself {
				sourceIP = localhost
			}
			reply.NewPayload(s.User.MAC, sourceIP, s.User.Username, sourceUsername)
			reply.Secret = stepMessage
			reply.ProtoType = protoType
			reply.Round = round + 1
			s.sendMessage(sourceIP, reply)
		}
		return
	default:
		// another type of error, which means err is probably not nil
		if err != nil {
			log.Panicf("ReceiveMessage: %s, Error Type: %s", err.Error(), errorType)
		}
	}
}

// Handles a chat message
func (s *Server) handleChatMessage(msg *ChatMessage) {
	sourceMAC, _, sourceUsername := msg.SourceID()
	messageYourself := sourceMAC == s.User.MAC && sourceUsername == s.User.Username
	sessions := s.GetSessionsWithFriend(sourceMAC, sourceUsername)
	friend := s.User.GetFriendByUsernameAndMAC(sourceUsername, sourceMAC)

	if len(sessions) == 0 || friend == nil {
		log.Panicln("No session or no friend from msg")
	}

	var sess *Session
	// There are two sessions, so grab the one that doesn't have the same timestamp as you
	if messageYourself {
		sess = sessions[1]
	} else {
		// There should only be one session between A -> B if you aren't messaging yourself, so grab that
		sess = sessions[0]
	}
	if db.GetSession(sess.Proto.GetSessionID()) == nil {
		sess.Save()
	}
	dec, _ := sess.Proto.Decrypt(msg.Text, s.onInfoReceive)
	if sess.Proto.IsActive() && dec[0] != nil {
		// Print the decoded message and IP
		currTime := time.Now()
		s.onReceiveChatMessage(dec[0], friend, currTime)
		db.InsertMessage(sess.Proto.GetSessionID(), dec[0], core.GetFormattedTime(currTime), db.Received)
	}
}

func initDialer(address string) (*websocket.Conn, error) {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(fmt.Sprintf("ws://%s/ws", address), http.Header{})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Start up server
func (s *Server) Start() error {
	// TODO wait until UI has started somehow?
	(*s).onInfoReceive("Launching server ...")

	fullAddr := fmt.Sprintf("%s:%d", localhost, Port)
	srv := &http.Server{Addr: fullAddr}
	http.HandleFunc("/ws", s.handleMessage)
	go srv.ListenAndServe()
	(*s).Listener = srv
	url, cmd, err := core.SetupTunnel(Port, (*s).User.Username, (*s).User.MAC)
	if err != nil {
		log.Panicln(err)
	}
	(*s).onInfoReceive(fmt.Sprintf("Your public url is: %s\n", url))
	(*s).Tunnel = cmd
	(*s).User.IP = url
	s.StartOTRSession(db.Self)
	return nil
}

// End server connection
func (s *Server) Shutdown() error {
	log.Println("Shutting down server...")
	if (*s).Tunnel != nil {
		if err := (*s).Tunnel.Process.Kill(); err != nil {
			return err
		}
		log.Println("Killed reverse-proxy tunnel")
	}
	if (*s).Listener != nil {
		return (*s).Listener.Close()
	}
	return nil
}

// Sends a formatted Message object with the server, after an active session between the two users have been established
func (s *Server) sendMessage(destIp string, msg Message) error {
	addr := fmt.Sprintf("%s:%d", destIp, Port)
	if destIp != localhost {
		addr = destIp
	}
	dialer, err := initDialer(addr)
	if err != nil {
		return err
	}
	w, err := dialer.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(w)
	if err = enc.Encode(&msg); err != nil {
		return err
	}
	return w.Close()
}

// Get all sessions that a user talks to an IP
// There are only 2 if a user is talking to himself
// otherwise only 1 session is returned
func (s *Server) GetSessionsWithFriend(friendMAC string, friendUsername string) []*Session {
	var filterSessions []*Session
	for _, sess := range *(*s).Sessions {
		if sess.To.MAC == friendMAC && sess.To.Username == friendUsername {
			filterSessions = append(filterSessions, sess)
		}
	}
	return filterSessions
}

// Start a session with a destination IP using a protocol
func (s *Server) StartOTRSession(displayName string) error {
	friend := s.User.GetFriendByDisplayName(displayName)
	if friend == nil {
		return errors.New(fmt.Sprintf("You do not have a friend named '%s'\n", displayName))
	}
	sessions := s.GetSessionsWithFriend(friend.MAC, friend.Username)
	if len(sessions) != 0 {
		return nil
	}

	proto := new(protocol.OTRProtocol)
	firstMessage, err := proto.NewSession()
	if err != nil {
		log.Printf("StartOTRSession: Error starting new session: %s", err)
		return err
	}

	// If messaging yourself, use your local IP as the sender too
	sourceIP := s.User.IP
	if displayName == db.Self {
		sourceIP = localhost
	}

	msg := new(HandshakeMessage)
	msg.NewPayload(s.User.MAC, sourceIP, s.User.Username, friend.Username)
	msg.Secret = []byte(firstMessage)
	msg.ProtoType = proto.ToType()
	msg.Round = 0

	err = s.sendMessage(friend.IP, msg)
	if err != nil {
		log.Printf("Error starting OTR session: %s\n", err.Error())
	}
	return err
}

// Sends a friend request to a specified destUsername@destIP
func (s *Server) SendFriendRequest(destIP, destUsername string) error {
	friendRequest := new(FriendMessage)
	friendRequest.NewPayload(s.User.MAC, s.User.IP, s.User.Username, destUsername)

	return s.sendMessage(destIP, friendRequest)
}

// Sends a chat message based on friend display name
func (s *Server) SendChatMessage(friendDisplayName, message string) error {
	chatMsg := new(ChatMessage)

	friend := s.User.GetFriendByDisplayName(friendDisplayName)
	if friend == nil {
		return errors.New(fmt.Sprintf("Friend with display name '%s' does not exist", friendDisplayName))
	}
	sessions := s.GetSessionsWithFriend(friend.MAC, friend.Username)
	if len(sessions) == 0 {
		return errors.New(fmt.Sprintf("Cannot communicate with '%s' without an active session\n", friendDisplayName))
	}
	userSession := sessions[0]
	cyp, err := userSession.Proto.Encrypt(chatMsg.Text)
	if err != nil {
		return err
	}
	(*chatMsg).Text = cyp[0]

	sourceIP := s.User.IP
	// If messaging yourself, use your local IP as the sender too
	if friend.IP == localhost {
		sourceIP = localhost
	}

	chatMsg.NewPayload(s.User.MAC, sourceIP, s.User.Username, friend.Username)
	bytes := []byte(message)
	(*chatMsg).Text = bytes

	if err := s.sendMessage(friend.IP, chatMsg); err != nil {
		// We had an issue with sending a message, so clear our session with the user
		userSession.EndSession()
	}
	// If we didn't have an issue, save the message into the database
	db.InsertMessage(userSession.Proto.GetSessionID(), bytes, core.GetFormattedTime(time.Now()), db.Sent)
	return nil
}
