package core

import (
	"bufio"
	"fmt"
	"github.com/wavyllama/chat/server"
	"os"
	"runtime"
	"strings"
)

const (
	exit         = ":exit"
	friend       = ":friend"
	display_name = ":display_name"
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
			if len(message) == 3 {
				// Todo: Make friend message
			} else {
				fmt.Printf("Format to add friend: '%s ipaddr username'\n", friend)
			}
		} else if len(words) == 1 {
			// Todo: Friend Exits
			activeFriend = words[0][1:]
		} else {
			fmt.Printf("Format for commands: '%s' or '%s'\n", exit, display_name)
		}
	} else {
		if activeFriend != "" {
			// Todo: Make chat message
		} else {
			fmt.Printf("Please set active friend: '%s'\n", display_name)
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
