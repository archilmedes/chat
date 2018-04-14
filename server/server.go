package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const (
	Port uint16 = 4242
	Network = "tcp"
)

// Simple Server struct
type Server struct {
	User *User
	Listener *net.TCPListener
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
func handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		log.Panicf("handleConnection: %s", err.Error())
	}
	fmt.Printf("%s: %s\n", msg.MAC, msg.Text)
}

// Function that continuously polls for new messages being sent to the server
func receive(listener *net.TCPListener) {
	for {
		if conn, err := (*listener).AcceptTCP(); err == nil {
			go handleConnection(conn)
		}
	}
}

// Start up server
func (s *Server) Start(username string, mac string, ip string) error {
	var err error
	log.Println("Launching Server...")
	(*s).User = &User{username, mac, ip}
	// TODO: Friends handling multi-cast and storing MAC address
	ipAddr := fmt.Sprintf("%s:%d", ip, Port)
	if (*s).Listener, err = setupServer(ipAddr); err != nil {
		return err
	}
	go receive((*s).Listener)
	log.Printf("Listening on: '%s:%d'", ip, Port)
	return nil
}

// End server connection
func (s *Server) Shutdown() error {
	log.Println("Shutting Down Server...")
	return (*s).Listener.Close()
}

func initDialer(address string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr(Network, address)
	if err != nil {
		return nil, err
	}
	return net.DialTCP(Network, nil, tcpAddr)
}

// Send a message to another Server
func (s *Server) Send(address string, MAC string, message []byte) error  {
	dialer, err := initDialer(fmt.Sprintf("%s:%d", address, Port))
	if err != nil {
		return err
	}
	var msg Message
	msg.Init(MAC, message)
	encoder := json.NewEncoder(dialer)
	if err = encoder.Encode(&msg); err != nil {
		return err
	}
	return nil
}
