package core

import (
	"bufio"
	"fmt"
	"github.com/wavyllama/chat/db"
	"os"
	"strings"
	"sync"
	"time"
)

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
	fmt.Printf("You have received a friend request from %s@%s (':accept' or ':reject:'):\n", username, ip)
	defer func() {
		Cond.Signal()
	}()
	Cond.L.Lock()
	Cond.L.Unlock()
	return getDisplayName()
}

// Gets the formatted input time to save in the database
func GetFormattedTime(t time.Time) string {
	timestampParts := strings.Split(t.String(), " ")
	return fmt.Sprintf("%s %s", timestampParts[0], timestampParts[1])
}
