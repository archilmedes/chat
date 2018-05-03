package main

import (
	"bufio"
	"fmt"
	"github.com/wavyllama/chat/core"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/wavyllama/chat/ui"
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/server"
)

func main() {
	db.SetupDatabase()
	if len(os.Args) == 2 && (os.Args[1] == "--reset" || os.Args[1] == "-r") {
		db.ClearDatabase()
	}

	mac, ip, err := core.GetAddresses()
	if err != nil {
		fmt.Printf("getAddresses: %s", err.Error())
	}

	user := core.Login(bufio.NewScanner(os.Stdin), ip)
	// Update the IP anyways
	user.IP = ip
	user.MAC = mac

	var program = server.InitServer()
	if err := program.Start(user); err != nil {
		log.Fatalf("main: %s", err.Error())
	}
	defer program.Shutdown()
	if _, err = ui.NewUI(program); err != nil {
		log.Fatalf("main: %s", err.Error())
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		os.Exit(0)
	}()
}
