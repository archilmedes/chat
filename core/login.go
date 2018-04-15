package core

import (
	"bufio"
	"regexp"
	"fmt"
	"strings"
	"os"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"chat/db"
)

const Me = "me"

var terminalReadPassword = terminal.ReadPassword

// Get username from stdin
func getUsername(scanner *bufio.Scanner) string {
	re := regexp.MustCompile("^[[:alnum:]]+$")
	for {
		fmt.Print("Username: ")
		scanner.Scan()
		username := strings.TrimSpace(scanner.Text())
		if strings.EqualFold(Me, username) {
			fmt.Printf("getUsername: %s is reserved!\n", username)
			continue
		}
		if re.MatchString(username) {
			return username
		}
		fmt.Printf("getUsername: %s is an invalid username!\n", username)
	}
	fmt.Println(scanner.Err().Error())
	os.Exit(1)
	return ""
}

// Returning user sign-in
func signIn(username string) bool {
	for counter := 0; counter < 3; counter++ {
		fmt.Print("Password: ")
		password, err := terminalReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if db.GetUser(username, string(password)) != nil {
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
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println()
		password := string(bytePassword)
		fmt.Print("Confirm password: ")
		bytePassword, err = terminalReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println()
		if password == string(bytePassword) {
			return db.AddUser(username, password, ip)
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
	if db.UserExists(username) {
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
