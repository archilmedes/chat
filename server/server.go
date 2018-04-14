package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"chat/db"
	"runtime"
	"chat/protocol"
)

const (
	Port uint16 = 4242
	Network = "tcp"
)

// Simple Server struct
type Server struct {
	User *db.User
	Listener *net.TCPListener
	Sessions []*Session
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
func handleConnection(s *Server, conn *net.TCPConn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		log.Panicf("handleConnection: %s", err.Error())
	}

	//fromUser := new(core.User)
	//fromUser.IP = conn.RemoteAddr().(*net.TCPAddr).IP.String()
	//dec, err := activeUser.ReceiveMessage(fromUser, msg.Text)
	//switch errorType := err.(type) {
	//default:
	//	log.Panicf("ReceiveMessage: %s, Error Type: %s", err.Error(), errorType)
	//case protocol.OTRHandshakeStep:
	//	// If it's part of the OTR handshake, send a message back directly, and return
	//	sendMessage(fromUser, dec)
	//	return
	//}
	//
	//fmt.Printf("%s: %s\n", msg.User, dec)
	if msg.StartProto != nil {
		fmt.Printf("They want to start a new protocol of type %s", msg.StartProto)
	}
	s.CreateOrGetSession(msg)
	fmt.Printf("%s: %s\n", msg.MAC, msg.Text)
}

// Function that continuously polls for new messages being sent to the server
func receive(s *Server) {
	for {
		if conn, err := (*(*s).Listener).AcceptTCP(); err == nil {
			go handleConnection(s, conn)
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
	go receive(s)
	log.Printf("Listening on: '%s:%d'", ip, Port)
	return nil
}

// End server connection
func (s *Server) Shutdown() error {
	log.Println("Shutting Down Server...")
	return (*s).Listener.Close()
}

// Send a message to another Server
func (s *Server) Send(address string, MAC string, message []byte) error  {
	dialer, err := initDialer(fmt.Sprintf("%s:%d", address, Port))
	if err != nil {
		return err
	}
	msg := NewMessage(MAC, address, message)
	encoder := json.NewEncoder(dialer)
	if err = encoder.Encode(&msg); err != nil {
		return err
	}
	return nil
}

func (s *Server) CreateOrGetSession(msg Message) (*Session) {
	for _, sess := range s.Sessions {
		if (*sess).ConverseWith(msg.IP) {
			return sess
		}
	}
	friend := new(User)
	friend.IP = msg.IP
	friend.MAC = msg.MAC
	return NewSession(s.User, friend, protocol.PlainProtocol{})
}
