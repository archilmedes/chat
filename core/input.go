package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var Friending FRFlag = DONE

func getDisplayName() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter display name: ")
	scanner.Scan()
	displayName := strings.TrimSpace(scanner.Text())
	if strings.ToLower(displayName) == Self {
		fmt.Println("Username is reserved! Please select another")
		getDisplayName()
	}
	return displayName
}

// Get display name from stdin
func GetDisplayNameFromConsole(ip string, username string) string {
	Friending = DONE
	fmt.Printf("You have received a friend request from %s@%s (':accept' or ':reject:'):", username, ip)
	defer func() {
		Friending = DONE
	}()
	for Friending == DONE {
	}
	if Friending == REJECT {
		return ""
	}
	return getDisplayName()
}
