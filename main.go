package main

import (
	"bufio"
	"fmt"
	"github.com/wavyllama/chat/protocol"
	"github.com/wavyllama/chat/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/wavyllama/chat/core"
	"github.com/wavyllama/chat/db"
)

func main() {
	db.SetupDatabase()
	mac, ip, err := core.GetAddresses()
	if err != nil {
		fmt.Printf("getAddresses: %s", err.Error())
	}
	username := core.Login(bufio.NewScanner(os.Stdin), ip)
	var program server.Server
	if err := program.Start("sameet", "b8:e8:56:2d:df:c2", "192.168.86.47"); err != nil {
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
	core.Listen(&program)
}
