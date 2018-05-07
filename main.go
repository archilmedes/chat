package main

import (
	"bufio"
	"github.com/wavyllama/chat/core"
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/server"
	"github.com/wavyllama/chat/ui"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db.SetupDatabase()
	if len(os.Args) == 2 && (os.Args[1] == "--reset" || os.Args[1] == "-r") {
		db.ClearDatabase()
		log.Println("Database cleared")
		db.SetupDatabase()
	}

	mac, ip, err := core.GetAddresses()
	if err != nil {
		log.Panicf("getAddresses: %s\n", err.Error())
	}

	user := core.Login(bufio.NewScanner(os.Stdin), ip)
	// Update the IP anyways
	user.IP = ip
	user.MAC = mac

	var program = server.InitServer(user)
	defer program.Shutdown()
	var uiCmpt *ui.UI
	if uiCmpt, err = ui.NewUI(program); err != nil {
		log.Fatalf("main: %s", err.Error())
	}

	if err = uiCmpt.Run(); err != nil {
		log.Fatalf("UI Run: %s\n", err.Error())
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		os.Exit(0)
	}()
}
