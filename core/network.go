package core

import (
	"errors"
	"net"
	"runtime"
	"fmt"
	"os/exec"
	"log"
)

// Get MAC and public IPv4 addresses
// Help from: https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
// Help from: http://grokbase.com/t/gg/golang-nuts/13cf1dcxhs/go-nuts-getting-ip-address-and-hardware-address-in-golang
func GetAddresses() (string, string, error) { // MAC address, IPv4 address, error
	var mac, ip string
	list, err := net.Interfaces()
	if err != nil {
		goto errOccurred
	}
	for _, iface := range list {
		if runtime.GOOS == "windows" && iface.Name != "Wi-Fi" {
			continue
		}
		if hardware := iface.HardwareAddr.String(); len(hardware) != 0 {
			mac = hardware
		}
		addrs, err := iface.Addrs()
		if err != nil {
			goto errOccurred
		}
		for _, addr := range addrs {
			if ipInfo, ok := addr.(*net.IPNet); ok && !ipInfo.IP.IsLoopback() {
				if ipInfo.IP.To4() != nil {
					ip = ipInfo.IP.String()
				}
			}
		}
	}
	if len(mac) != 0 && len(ip) != 0 {
		return mac, ip, nil
	}
errOccurred:
	return mac, ip, errors.New("cannot find MAC and IPv4 addresses")
}

// Set up an encrypted tunnel on a port and return the Cmd through the channel
func SetupTunnel(username string, port uint16, c chan *exec.Cmd) {
	args := fmt.Sprintf("--port %d --subdomain %s", port, username)
	cmd := exec.Command("lt", args)
	err := cmd.Start()
	if err != nil {
		log.Fatalln(err.Error())
	}
	c <- cmd
}
