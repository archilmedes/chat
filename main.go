package main

import (
	"bufio"
	"fmt"
	"github.com/wavyllama/chat/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/wavyllama/chat/core"
	"github.com/wavyllama/chat/db"
	"strings"
	"runtime"
	"github.com/wavyllama/chat/protocol"
)

const (
	exit        = ":exit"
	friend      = ":friend"
	displayName = ":displayName"
)

var activeFriend = ""

func handleInput(program *server.Server, message string) {
	words := strings.Fields(message)
	if len(message) == 0 {
		return
	}
	if message[0] == ':' {
		if words[0] == exit {
			runtime.Goexit()
		} else if words[0] == friend {
			if len(words) == 2 {
				friendInfo := strings.Split(words[1], "@")
				if len(friendInfo) == 2 {
					err := program.SendFriendRequest(friendInfo[1], friendInfo[0])
					if err != nil {
						log.Printf("Error sending friend request: %s\n", err)
					}
					return
				}
			}
			fmt.Printf("Format to add friend: '%s username@ipaddr'\n", friend)
		} else if len(words) == 1 && program.User.IsFriendsWith(words[0][1:]) {
			activeFriend = words[0][1:]
			program.StartSession(activeFriend, protocol.OTRProtocol{})
		} else {
			fmt.Printf("Format for commands: '%s' or '%s'\n", exit, displayName)
		}
	} else {
		if activeFriend != "" {
			err := program.SendChatMessage(activeFriend, message)
			if err != nil {
				log.Printf("Error sending message: %s\n", err)
			}
		} else {
			fmt.Printf("Please set active friend: '%s'\n", displayName)
		}
	}
}

// Listen to stdin for messages to be sent
func Listen(program *server.Server) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		handleInput(program, scanner.Text())
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
	Listen(&program)
}
