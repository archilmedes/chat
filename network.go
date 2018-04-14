package main

import (
	"net"
	"errors"
	"runtime"
)

// Get public IPv4 and MAC address
// Help from: https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
// Help from: http://grokbase.com/t/gg/golang-nuts/13cf1dcxhs/go-nuts-getting-ip-address-and-hardware-address-in-golang
func getAddresses() (string, string, error) { // MAC address, IPv4 address, error
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
