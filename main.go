package main

import (
	"bufio"
	"fmt"
	"github.com/wavyllama/chat/core"
	"github.com/wavyllama/chat/server"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"github.com/wavyllama/chat/db"
)

const (
	exit        = ":exit"
	friend      = ":friend"
	accept      = ":accept"
	reject      = ":reject"
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
		} else if words[0] == accept {
			core.Friending = core.ACCEPT
			core.Cond.Signal()
			core.Cond.L.Lock()
			for core.Friending == core.ACCEPT {
				core.Cond.Wait()
			}
			core.Cond.L.Unlock()
		} else if words[0] == reject {
			core.Friending = core.REJECT
			core.Cond.Signal()
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
			program.StartOTRSession(activeFriend)
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

	// TODO would be good to change this to return a User object, we could use this to load conversation history, etc.
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
