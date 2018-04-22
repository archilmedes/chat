package server

import (
	"errors"
	"fmt"
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/protocol"
	"log"
	"net"
	"time"
	"encoding/gob"
	"encoding/json"
)

const (
	Port    uint16 = 4242
	Network        = "tcp"
)

// Server holds the user and all of his sessions
type Server struct {
	User     *db.User
	Listener *net.TCPListener
	Sessions *[]Session
}

func init() {
	gob.Register(&FriendMessage{})
	gob.Register(&HandshakeMessage{})
	gob.Register(&ChatMessage{})
}

// Setup listener for the server
func setupServer(address string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr(Network, address)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP(Network, tcpAddr)
}

// Handle receiving messages from a TCPConn
func (s *Server) handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)
	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		log.Panicf("Error decoding message: %s", err.Error())
	}
	res, _ := json.Marshal(msg)
	fmt.Printf("RECEIVED MESSAGE: %s\n", string(res))

	sourceIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	sourceMAC, sourceUsername := msg.SourceID()
	if sourceMAC == "" || sourceUsername == "" {
		log.Panicln("Received ill-formatted message")
	}
	if sourceUsername != s.User.Username {
		fmt.Println("Received a message but it was not for me.")
		return
	}
	messageYourself := sourceMAC == s.User.MAC && sourceUsername == s.User.Username
	sessions := s.GetSessionsWithFriend(sourceMAC, sourceUsername)
	friend := s.User.GetFriend(sourceUsername, sourceMAC)

	switch msg.(type) {
	case *FriendMessage:
		// TODO add friend if not necessary
		fmt.Println("Received a friend message")
	case *HandshakeMessage:
		// We are in a handshake, so the friend should exist already
		// TODO uncomment when friend request is implemented
		//if friend == nil {
		//	log.Panicln("You must be a friend to participate in a handshake")
		//}
		var createdSession bool
		var sess Session
		round := msg.(*HandshakeMessage).Round
		protoType, startSessionTime := msg.(*HandshakeMessage).ProtoType, msg.(*HandshakeMessage).SessionTime

		// In a handshake, create a new session if there aren't the required number of sessions in either situation
		if len(sessions) != 2 && messageYourself || (len(sessions) != 1 && !messageYourself) {
			// TODO remove when friend request is established
			if friend == nil {
				friend = new(db.Friend)
				friend.Username = sourceUsername
				friend.MAC = sourceMAC
			}
			sess = *NewSessionFromUserAndMessage(s.User, friend, protoType, startSessionTime)
			*(*s).Sessions = append(*(*s).Sessions, sess)
			createdSession = true
		} else if len(sessions) == 2 && messageYourself {
			// Communicating between yourself, rotate sessions based on round (even/odd)
			sess = sessions[round % 2]
		} else {
			sess = sessions[0]
		}

		dec, err := sess.Proto.Decrypt(msg.(*HandshakeMessage).Secret)

		switch errorType := err.(type) {
		case protocol.OTRHandshakeStep:
			// Send each part of the handshake message back and immediately return
			for _, stepMessage := range dec {
				reply := new(HandshakeMessage)
				reply.NewPayload(s.User.MAC, s.User.Username, sourceUsername)
				reply.Secret = stepMessage
				reply.ProtoType = msg.(*HandshakeMessage).ProtoType
				// If we created a session here, then set current time as start time
				if createdSession {
					reply.SessionTime = time.Now()
				}
				reply.Round = round + 1
				s.sendMessage(friend.DisplayName, sourceIP, reply)
			}
			return
		default:
			// another type of error, which means err is probably not nil
			if err != nil {
				log.Panicf("ReceiveMessage: %s, Error Type: %s", err.Error(), errorType)
			}
		}
	case *ChatMessage:
		var sess Session
		// There are two sessions, so grab the one that doesn't have the same timestamp as you
		if messageYourself {
			sess = sessions[1]
		} else {
			// There should only be one session between A -> B if you aren't messaging yourself, so grab that
			sess = sessions[0]
		}
		dec, _ := sess.Proto.Decrypt(msg.(*ChatMessage).Text)
		if sess.Proto.IsActive() && dec[0] != nil {
			// Print the decoded message and IP
			fmt.Printf("%s: %s\n", friend.DisplayName, dec[0])
		}
	}
}

// Function that continuously polls for new messages being sent to the server
func (s *Server) receive() {
	for {
		if conn, err := (*(*s).Listener).AcceptTCP(); err == nil {
			go s.handleConnection(conn)
		}
	}
}

func initDialer(address string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr(Network, address)
	if err != nil {
		return nil, err
	}
	return net.DialTCP(Network, nil, tcpAddr)
}

// Start up server
func (s *Server) Start(username string, mac string, ip string) error {
	var err error
	log.Println("Launching Server...")
	(*s).User = &db.User{username, mac, ip}
	ipAddr := fmt.Sprintf("%s:%d", ip, Port)
	if (*s).Listener, err = setupServer(ipAddr); err != nil {
		return err
	}
	// Initialize the session struct to a pointer
	(*s).Sessions = &[]Session{}
	go s.receive()
	log.Printf("Listening on: '%s:%d'", ip, Port)
	return nil
}

// End server connection
func (s *Server) Shutdown() error {
	log.Println("Shutting Down Server...")
	return (*s).Listener.Close()
}

// Sends a formatted Message object with the server, after an active session between the two users have been established
func (s *Server) sendMessage(destIp string, msg Message) error {
	dialer, err := initDialer(fmt.Sprintf("%s:%d", destIp, Port))
	if err != nil {
		return err
	}

	// Unless you're handshaking, then you must have an active session to send a message
	chatMsg, ok := msg.(*ChatMessage)
	if ok {
		friend := s.User.GetFriend()
		sessions := s.GetSessionsWithFriend(friend.MAC, msg.DestID())
		if len(sessions) == 0 {
			return errors.New(fmt.Sprintf("Cannot communicate with %s without an active session\n", msg.DestID()))
		}
		cyp, err := sessions[0].Proto.Encrypt(chatMsg.Text)
		if err != nil {
			return err
		}
		(*chatMsg).Text = cyp[0]
	}
	res, _ := json.Marshal(msg)
	fmt.Printf("Sending message %s\n", string(res))
	encoder := gob.NewEncoder(dialer)
	if err := encoder.Encode(&msg); err != nil {
		return err
	}
	return nil
}

// Send a message to another Server
func (s *Server) Send(destUsername, displayName, destIp string, message []byte) error {
	chatMsg := new(ChatMessage)
	chatMsg.NewPayload(s.User.MAC, s.User.IP, destUsername)
	chatMsg.Text = message
	return s.sendMessage(destIp, chatMsg)
}

// Get all sessions that a user talks to an IP
// There are only 2 if a user is talking to himself
// otherwise only 1 session is returned
func (s *Server) GetSessionsWithFriend(friendMAC string, friendUsername string) []Session {
	var filterSessions []Session
	for _, sess := range *(*s).Sessions {
		if sess.To.MAC == friendMAC && sess.To.Username == friendUsername {
			filterSessions = append(filterSessions, sess)
		}
	}
	return filterSessions
}

// Start a session with a destination IP using a protocol
func (s *Server) StartSession(destUsername, destIp string, proto protocol.Protocol) error {
	firstMessage, err := proto.NewSession()
	if err != nil {
		log.Panicf("StartSession: Error starting new session: %s", err)
		return err
	}

	msg := new(HandshakeMessage)
	msg.NewPayload(s.User.MAC, s.User.Username, destUsername)
	msg.Secret = []byte(firstMessage)
	msg.ProtoType = proto.ToType()
	msg.Round = 0
	return s.sendMessage(destIp, msg)
}
