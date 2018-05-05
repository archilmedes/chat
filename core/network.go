package core

import (
	"errors"
	"net"
	"runtime"
	"os/exec"
	"bufio"
	"strings"
	"strconv"
	"fmt"
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

func SetupTunnel(port uint16, username, mac string) (string, *exec.Cmd, error) {
	macNoComma := strings.Replace(mac, ":", "", -1)
	subDomain := fmt.Sprintf("%s-%s", username, macNoComma)
	cmd := exec.Command("lt", "--port", strconv.Itoa(int(port)), "--subdomain", subDomain)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", cmd, err
	}

	if err := cmd.Start(); err != nil {
		return "", cmd, err
	}

	output := bufio.NewScanner(stdout)
	var url string
	for output.Scan() {
		url = strings.Split(output.Text(), "your url is: ")[1]
		return url, cmd, output.Err()
	}
	return "", cmd, errors.New("could not run the command")
}
