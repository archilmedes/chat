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

const Exit = "exit"

func main() {
	var program server.Server
	if err := program.Start("Archil"); err != nil {
		log.Fatalf("main: %s", err.Error())
	}
	defer program.Shutdown()
	sig := make (chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		os.Exit(0)
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if message == Exit {
			os.Exit(0)
		}
		stringSlice := strings.Fields(message)
		if err := program.Send(stringSlice[0], strings.Join(stringSlice[1:], " ")); err != nil {
			fmt.Printf("input: %s", err.Error())
		}
	}
	if scanner.Err() != nil {
		fmt.Printf(scanner.Err().Error())
		os.Exit(1)
	}
}