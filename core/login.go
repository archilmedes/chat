// The core package provides fundamental core functionality shared by multiple parts
// of the program
package core

import (
	"bufio"
	"fmt"
	"github.com/wavyllama/chat/db"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"regexp"
	"strings"
	"syscall"
)

const Self = "me" // Display name for self.

var dbAddUser = db.AddUser
var dbGetUser = db.GetUser
var dbUserExists = db.UserExists
var terminalReadPassword = terminal.ReadPassword

// Get username from stdin
func getUsername(scanner *bufio.Scanner) string {
	re := regexp.MustCompile("^[[:alnum:]]+$")
	for {
		fmt.Print("Username: ")
		scanner.Scan()
		username := strings.TrimSpace(scanner.Text())
		if strings.EqualFold(Self, username) {
			fmt.Printf("getUsername: %s is reserved!\n", username)
			continue
		}
		if re.MatchString(username) {
			return username
		}
		fmt.Printf("getUsername: %s is an invalid username!\n", username)
	}
	log.Fatalln(scanner.Err().Error())
	return ""
}

// Sign-in for returning user
func signIn(username string) bool {
	for counter := 0; counter < 3; counter++ {
		fmt.Print("Password: ")
		password, err := terminalReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if dbGetUser(username, string(password)) != nil {
			return true
		}
		fmt.Printf("signIn: invalid password!\n")
	}
	return false
}

// Create an account for a new user
func createAccount(username string, ip string) bool {
	for counter := 0; counter < 3; counter++ {
		fmt.Print("Enter new password: ")
		bytePassword, err := terminalReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		password := string(bytePassword)
		fmt.Print("Confirm password: ")
		bytePassword, err = terminalReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if password == string(bytePassword) {
			return dbAddUser(username, password, ip)
		} else {
			fmt.Println("Passwords do not match!")
		}
	}
	return false
}

// Login user
func Login(scanner *bufio.Scanner, ip string) string {
	username := getUsername(scanner)
	var successful bool
	if dbUserExists(username) {
		successful = signIn(username)
	} else {
		successful = createAccount(username, ip)
	}
	if successful {
		return username
	} else {
		return Login(scanner, ip)
	}
}
