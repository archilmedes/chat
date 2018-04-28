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

const (
	numTries = 3
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
func signIn(username string) *db.User {
	for counter := 0; counter < numTries; counter++ {
		fmt.Print("Password: ")
		password, err := terminalReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		user := dbGetUser(username, string(password))
		if user != nil {
			return user
		}
		fmt.Println("Username and password combination do not match.")
	}
	fmt.Printf("Authentication failed after %d attempts.\n", numTries)
	return nil
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
func createAccount(username string) *db.User {
	for counter := 0; counter < numTries; counter++ {
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
		if password == string(bytePassword) && dbAddUser(username, password, "doesntmatter") {
			return dbGetUser(username, password)
		} else {
			fmt.Println("Passwords do not match!")
		}
	}
	fmt.Printf("Account creation failed after %d attempts.\n", numTries)
	return nil
}

// Login user
func Login(scanner *bufio.Scanner) *db.User {
	username := getUsername(scanner)
	var user *db.User
	if dbUserExists(username) {
		user = signIn(username)
	} else {
		user = createAccount(username)
	}
	if user != nil{
		return user
	} else {
		return Login(scanner)
	}
}
