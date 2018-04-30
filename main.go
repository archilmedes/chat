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
	"runtime"
	"strings"
	"syscall"
)

const (
	exit        = ":exit"
	friend      = ":friend"
	accept      = ":accept"
	reject      = ":reject"
	displayName = ":displayName"
	unfriend    = ":delete"
	deleteSelf  = ":deleteSelf"
)

var activeFriend = ""

func handleSpecialString(program *server.Server, words []string) {
	switch words[0] {
	case exit:
		runtime.Goexit()
	case accept:
		core.Friending = core.ACCEPT
		core.Cond.Signal()
		core.Cond.L.Lock()
		for core.Friending == core.ACCEPT {
			core.Cond.Wait()
		}
		core.Cond.L.Unlock()
	case reject:
		core.Friending = core.REJECT
		core.Cond.Signal()
	case friend:
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
	case unfriend:
		if len(words) == 2 {
			if !program.User.DeleteFriend(words[1]) {
				fmt.Printf("Error deleting friend: %s\n", words[1])
			}
		}
		fmt.Printf("Format to delete a friend: '%s displayName\n", unfriend)
	case deleteSelf:
		if !program.User.Delete() {
			fmt.Printf("Failed to delete your account\n")
		}
		fmt.Println("Successfully deleted all of your data")
		os.Exit(0)
	default:
		if len(words) == 1 && program.User.IsFriendsWith(words[0][1:]) {
			activeFriend = words[0][1:]
			program.StartOTRSession(activeFriend)
		} else {
			fmt.Printf("Format for commands: '%s' or '%s'\n", exit, displayName)
		}
	}
}

func handleInput(program *server.Server, message string) {
	words := strings.Fields(message)
	if len(message) == 0 {
		return
	}
	if message[0] == ':' {
		handleSpecialString(program, words)
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

	user := core.Login(bufio.NewScanner(os.Stdin), ip)
	// Update the IP anyways
	user.IP = ip
	user.MAC = mac
	// Update the IP in the database
	user.UpdateMyIP()
	var program server.Server
	if err := program.Start(user); err != nil {
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
