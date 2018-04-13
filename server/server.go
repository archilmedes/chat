package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
)

const (
	Port uint16 = 4242
	Network = "tcp"
)

// Simple Server struct
type Server struct {
	IP, Username string
	Listener *net.TCPListener
}

// Get public ip address
// Help from: https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func getIp() (string, error) {
	list, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range list {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			if ipInfo, ok := addr.(*net.IPNet); ok && !ipInfo.IP.IsLoopback() {
				if ((runtime.GOOS == "windows" && iface.Name == "Wi-Fi") ||
					runtime.GOOS != "windows") && ipInfo.IP.To4() != nil {
					return ipInfo.IP.String(), nil
				}
			}
		}
	}
	return "", errors.New("getIp: cannot find public ip address")
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
	fmt.Printf("%s: %s\n", msg.User, msg.Text)
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
func (s *Server) Start(username string) error {
	log.Println("Launching Server...")
	(*s).Username = username
	var err error
	if (*s).IP, err = getIp(); err != nil {
		return err
	}
	ipAddr := fmt.Sprintf("%s:%d", (*s).IP, Port)
	if (*s).Listener, err = setupServer(ipAddr); err != nil {
		return err
	}
	go receive((*s).Listener)
	log.Printf("Listening on: '%s:%d'", (*s).IP, Port)
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
func (s *Server) Send(address string, message string) error  {
	dialer, err := initDialer(fmt.Sprintf("%s:%d", address, Port))
	if err != nil {
		return err
	}
	var msg Message
	msg.Init((*s).Username, message)
	encoder := json.NewEncoder(dialer)
	if err = encoder.Encode(&msg); err != nil {
		return err
	}
	return nil
}
