package core

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"github.com/wavyllama/chat/db"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"testing"
)

func setupTest(file string) (*os.File, *bufio.Scanner) {
	os.Stdin.Close()
	f, _ := os.Open(file)
	scanner := bufio.NewScanner(f)
	dbAddUser = func(username string, password string, ipAddress string) bool {
		return true
	}
	dbGetUser = func(username string, password string) *db.User {
		return new(db.User)
	}
	dbUserExists = func(username string) bool {
		return username == "bob"
	}
	terminalReadPassword = func(fd int) ([]byte, error) {
		scanner.Scan()
		return scanner.Bytes(), nil
	}
	return f, scanner
}

func TestLogin_New_User(t *testing.T) {
	f, scanner := setupTest("login_test_new_user.txt")
	defer f.Close()
	defer func() {
		dbAddUser = db.AddUser
		dbGetUser = db.GetUser
		dbUserExists = db.UserExists
		terminalReadPassword = terminal.ReadPassword
	}()
	assert.Equal(t, Login(scanner, ""), "sameertqa")
}

func TestLogin_Current_User(t *testing.T) {
	f, scanner := setupTest("login_test_current_user.txt")
	defer f.Close()
	defer func() {
		dbAddUser = db.AddUser
		dbGetUser = db.GetUser
		dbUserExists = db.UserExists
		terminalReadPassword = terminal.ReadPassword
	}()
	assert.Equal(t, Login(scanner, ""), "bob")
}
