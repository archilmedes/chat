package main

import (
	"bufio"
	"fmt"
	"github.com/wavyllama/chat/core"
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db.SetupDatabase()
	mac, ip, err := core.GetAddresses()
	if err != nil {
		fmt.Printf("getAddresses: %s", err.Error())
	}

	user := core.Login(bufio.NewScanner(os.Stdin), ip)
	// Update the IP anyways
	user.IP = ip
	user.MAC = mac
	// Update the IP in the database
	user.UpdateMyIP()
	var program = new(server.Server)
	if err := program.Start(user); err != nil {
		log.Fatalf("main: %s", err.Error())
	}
	defer program.Shutdown()
	if _, err = server.NewUI(program); err != nil {
		log.Fatalf("main: %s", err.Error())
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		os.Exit(0)
	}()
}
