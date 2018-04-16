package main

import (
	"bufio"
	"fmt"
	"github.com/wavyllama/chat/core"
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/protocol"
	"github.com/wavyllama/chat/server"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const Exit = "exit"

// Listen to standard in for messages to be sent
func listen(program *server.Server) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if message == Exit {
			return
		}
		stringSlice := strings.Fields(message)
		// Message format is: "IP message"
		if err := program.Send(stringSlice[0], strings.Join(stringSlice[1:], " ")); err != nil {
			fmt.Printf("input: %s\n", err.Error())
		}
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err().Error())
		return
	}
}

func main() {
	db.SetupDatabase()
	mac, ip, err := core.GetAddresses()
	if err != nil {
		fmt.Printf("getAddresses: %s", err.Error())
	}
	username := core.Login(bufio.NewScanner(os.Stdin), ip)
	var program server.Server
	if err := program.Start(username, mac, ip); err != nil {
		log.Fatalf("main: %s", err.Error())
	}
	defer program.Shutdown()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		os.Exit(0)
	}()
	program.StartSession(ip, protocol.OTRProtocol{})
	listen(&program)
}
