package core

import (
	"fmt"
	"strings"
	"os"
	"bufio"
)

// Get display name from stdin
func GetDisplayNameFromConsole() string {
	var username string
	for {
		fmt.Print("Enter display name: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		username = strings.TrimSpace(scanner.Text())
		if strings.ToLower(username) == Self {
			fmt.Println("Username is reserved! Please select another")
			continue
		}
		break
	}
	return username
}
