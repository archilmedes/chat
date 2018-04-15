package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"chat/protocol"
	"time"
	"errors"
	"chat/db"
)

const (
	Port uint16 = 4242
	Network = "tcp"
)

// Simple Server struct
type Server struct {
	User *db.User
	Listener *net.TCPListener
	Sessions *[]Session
}

// Setup listener for server
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
	decoder := json.NewDecoder(conn)
	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		log.Panicf("handleConnection: %s", err.Error())
	}

	if s.User.IP != msg.DestIP {
		log.Panicln("User received a message that was not meant for them")
	}

	var sess Session
	oldNum := len(*(*s).Sessions)

	// If part of the handshake
	if msg.Handshake {
		idx := msg.ID % 2
		sessions := s.GetSessionsToIP(msg.DestIP) // One for normal cases, two for communicating to yourself
		if len(sessions) != 2 {
			// sess = s.CreateOrGetSession(msg)

			// If sending message to yourself TODO: handle other cases later
			if msg.SourceIP == msg.DestIP {
				sess = *NewSessionFromUserAndMessage(s.User, msg)
				*(*s).Sessions = append(*(*s).Sessions, sess)
			}
		} else {
			sess = sessions[idx]
		}
	} else {
		sess = s.CreateOrGetSession(msg)
	}
	newNum := len(*(*s).Sessions)
	createdSession := oldNum != newNum


	dec, err := sess.Proto.Decrypt([]byte(msg.Text))
	fmt.Println(string(dec[0]))

	switch errorType := err.(type) {
	case protocol.OTRHandshakeStep:
		// If it's part of the OTR handshake, send each part of the message back directly to the source,
		// and immediately return

		for _, stepMessage := range dec {
			reply := NewMessage(s.User, msg.SourceIP, string(stepMessage))
			reply.StartProtocol(sess.Proto)
			if createdSession {
				reply.StartProtoTimestamp = time.Now()
			}
			reply.ID = msg.ID + 1
			s.sendMessage(reply)
		}
		return
	default:
		if err != nil {
			log.Panicf("ReceiveMessage: %s, Error Type: %s", err.Error(), errorType)
		}
		if createdSession { // That means msg.Handshake must be true
			reply := NewMessage(s.User, msg.SourceIP, msg.Text)
			reply.StartProtocol(sess.Proto)
			if createdSession {
				reply.StartProtoTimestamp = time.Now()
			}
			reply.ID = msg.ID + 1
			s.sendMessage(reply)
		}
	}
	if sess.Proto.IsActive() && dec[0] != nil {
		// Print the decoded message and IP
		fmt.Printf("%s: %s\n", msg.SourceIP, dec[0])
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

// Private-helper method that sends a formatted message object with the server
func (s *Server) sendMessage(msg *Message) error {
	dialer, err := initDialer(fmt.Sprintf("%s:%d", msg.SourceIP, Port))
	if err != nil {
		return err
	}

	// Unless you're handshaking, then you must have an active session to send a message
	sessions := s.GetSessionsToIP((*msg).DestIP)
	if len(sessions) == 0 && !msg.Handshake {
		return errors.New(fmt.Sprintf("Cannot communicate with %s without an active session\n", msg.DestIP))
	} else if len(sessions) != 0 && !msg.Handshake {
		(*msg).StartProtoTimestamp = sessions[0].StartTime
		cyp, err := sessions[0].Proto.Encrypt([]byte((*msg).Text))
		if err != nil {
			return err
		}
		(*msg).Text = string(cyp[0])
	}

	encoder := json.NewEncoder(dialer)
	if err := encoder.Encode(msg); err != nil {
		return err
	}
	return nil
}

// Send a message to another Server
func (s *Server) Send(destIp string, message string) error  {
	return s.sendMessage(NewMessage(s.User, destIp, message))
}

// Get all sessions that a user talks to an IP. There are only 2 if a user is talking to himself
func (s *Server) GetSessionsToIP(ip string) []Session {
	var filterSessions []Session
	for _, sess := range *(*s).Sessions {
		if sess.ConverseWith(ip) {
			filterSessions = append(filterSessions, sess)
		}
	}
	return filterSessions
}

// Creates or gets a session based on the message received
func (s *Server) CreateOrGetSession(msg Message) Session {
	for _, sess := range *(*s).Sessions {
		if sess.ConverseWith(msg.SourceIP) && sess.StartTime != msg.StartProtoTimestamp {
			return sess
		}
	}
	sess := NewSessionFromUserAndMessage(s.User, msg)
	*(*s).Sessions = append(*(*s).Sessions, *sess)
	return *sess
}

// Start a session with a destination IP using a protocol
func (s *Server) StartSession(destIp string, proto protocol.Protocol) (error) {
	firstMessage, err := proto.NewSession()
	if err != nil {
		log.Panicf("StartSession: Error starting new session: %s", err)
		return err
	}

	msg := NewMessage(s.User, destIp, firstMessage)
	msg.StartProtocol(proto)
	msg.ID = 0
	return s.sendMessage(msg)
}
