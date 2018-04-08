package main

import (
	"bufio"
	"chat/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const Exit = "exit\n"

func main() {
	var program server.Server
	if err := program.Start("Archil"); err != nil {
		log.Fatalf("main: %s", err.Error())
	}
	sig := make (chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		program.Shutdown()
		os.Exit(0)
	}()
	for {
		var message string
		var err error
		reader := bufio.NewReader(os.Stdin)
		message, err = reader.ReadString('\n')
		if err != nil {
			program.Shutdown()
			log.Fatalf("main: %s", err)
			os.Exit(1)
		}
		if message == Exit {
			program.Shutdown()
			os.Exit(0)
		}
		stringSlice := strings.Fields(message)
		if err = program.Send(stringSlice[0], strings.Join(stringSlice[1:], " ")); err != nil {
			fmt.Errorf("main: %s", err.Error())
		} else {
			fmt.Println("Try again!")
		}
	}
}