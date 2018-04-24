package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"github.com/wavyllama/chat/db"
)

var Friending FRFlag = DONE
var mutex = sync.Mutex{}
var Cond = sync.NewCond(&mutex)
var CondWait = Cond.Wait
var bufioNewScanner = bufio.NewScanner

func getDisplayName() string {
	scanner := bufioNewScanner(os.Stdin)
	fmt.Println("Enter display name: ")
	scanner.Scan()
	displayName := strings.TrimSpace(scanner.Text())
	if strings.ToLower(displayName) == db.Self {
		fmt.Println("Username is reserved! Please select another.")
		displayName = getDisplayName()
	}
	return displayName
}

// Get display name from stdin
func GetDisplayNameFromConsole(ip string, username string) string {
	Friending = DONE
	fmt.Printf("You have received a friend request from %s@%s (':accept' or ':reject:'):\n", username, ip)
	defer func() {
		Friending = DONE
		Cond.Signal()
	}()
	Cond.L.Lock()
	for Friending == DONE {
		CondWait()
	}
	Cond.L.Unlock()
	if Friending == REJECT {
		return ""
	}
	return getDisplayName()
}
