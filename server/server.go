package server

import (
	"errors"
	"fmt"
	"net"
)

const (
	Port uint16  = 4242
	Network = "tcp"
)

type Server interface {
	Start() error
	Shutdown() error
	Send()
	Recv()
}

type SimpleServer struct {
	IP string
	Listener *net.TCPListener
}

/**
 * Get public ip address.
 * Help from: https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
 */
func getIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, i := range addrs {
		if ipInfo, ok := i.(*net.IPNet); ok && !ipInfo.IP.IsLoopback() {
			if ipInfo.IP.To4() != nil {
				return ipInfo.IP.String(), nil
			}
		}
	}
	return "", errors.New("getIp: cannot find public ip address")
}

/**
 * Setup listener for server.
 */
func setupServer() (*net.TCPListener, error) {
	address := fmt.Sprintf("localhost:%d", Port)
	tcpAddr, err := net.ResolveTCPAddr(Network, address)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP(Network, tcpAddr)
}

/**
 * Start up server.
 */
func (s *SimpleServer) Start() error {
	var err error
	fmt.Println("Launching Server...")
	(*s).IP, err = getIp()
	if err != nil {
		return err
	}
	(*s).Listener, err = setupServer()
	if err != nil {
		return err
	}
	fmt.Printf("Listening on: '%s:%d'", (*s).IP, Port)
	return nil
}

/**
 * End server connection.
 */
func (s *SimpleServer) Shutdown() error {
	return (*s).Listener.Close()
}

/**
 * Send
 */
