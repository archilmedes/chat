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
	"unicode"
)

var dbAddUser = db.AddUser
var dbGetUser = db.GetUser
var dbUserExists = db.UserExists
var terminalReadPassword = terminal.ReadPassword

// Get username from stdin
func getUsername(scanner *bufio.Scanner) string {
	re := regexp.MustCompile("^[[:alnum:]]{2,16}$")
	for {
		fmt.Print("Username: ")
		scanner.Scan()
		username := strings.TrimSpace(scanner.Text())
		if strings.EqualFold(db.Self, username) {
			fmt.Printf("getUsername: %s is reserved!\n", username)
			continue
		}
		if re.MatchString(username) {
			return username
		}
		fmt.Printf("%s is an invalid username!\n", username)
		fmt.Println("Usernames only contain alphanumeric characters only of max length 16")
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

func verifyPassword(try string) bool {
	lower, upper, num := false, false, false
	size := len(try)
	for _, c := range try {
		if unicode.IsNumber(c) {
			num = true
		} else if unicode.IsLower(c) {
			lower = true
		} else if unicode.IsUpper(c) {
			upper = true
		} else if !unicode.IsPrint(c) || unicode.IsSpace(c) {
			return false
		}
	}
	return lower && upper && num && 8 <= size && size <= 32
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
		if !verifyPassword(password) || strings.ToLower(password) == strings.ToLower(username) {
			fmt.Println("Invalid Password! Password must be between 8-32 characters long.")
			fmt.Println("Password must consist of at least one number and one uppercase and one lowercase character.")
			continue
		}
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
