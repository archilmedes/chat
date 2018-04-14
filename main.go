package main

import (
	"bufio"
	"chat/server"
	"chat/core"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"regexp"
	"chat/db"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	Exit = "exit"
	Me = "me"
)

// Listen to standard in for messages to be sent
func listen(program *server.Server) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if message == Exit {
			return
		}
		stringSlice := strings.Fields(message)
		if err := program.Send(stringSlice[0], "", []byte(strings.Join(stringSlice[1:], " "))); err != nil {
			fmt.Printf("input: %s\n", err.Error())
		}
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err().Error())
		return
	}
}

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
		password, err := terminal.ReadPassword(int(syscall.Stdin))
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
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println()
		password := string(bytePassword)
		fmt.Print("Confirm password: ")
		bytePassword, err = terminal.ReadPassword(int(syscall.Stdin))
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

func login(scanner *bufio.Scanner, ip string) string {
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
		return login(scanner, ip)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	db.SetupDatabase()
	mac, ip, err := core.GetAddresses()
	if err != nil {
		fmt.Printf("getAddresses: %s", err.Error())
	}
	username := login(scanner, ip)
	var program server.Server
	if err := program.Start(username, mac, ip); err != nil {
		log.Fatalf("main: %s", err.Error())
	}
	defer program.Shutdown()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		os.Exit(0)
	}()
	listen(&program)
}
