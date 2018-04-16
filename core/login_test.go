package core

import (
	"os"
	"testing"
	"chat/db"
	"github.com/stretchr/testify/assert"
	"bufio"
	"golang.org/x/crypto/ssh/terminal"
)

func setup() {
	db.SetupTestDatabaseAltDir()
	os.Stdin.Close()
}

func TestLogin_New_User(t *testing.T) {
	setup()
	f, _ := os.Open("login_test_new_user.txt")
	defer f.Close()
	defer func() {
		terminalReadPassword = terminal.ReadPassword
	}()
	scanner := bufio.NewScanner(f)
	terminalReadPassword = func (fd int) ([]byte, error) {
		scanner.Scan()
		return scanner.Bytes(), nil
	}
	assert.Equal(t, Login(scanner, ""), "sameertqa")
	db.DeleteUser("sameertqa")
}

func TestLogin_Current_User(t *testing.T) {
	setup()
	f, _ := os.Open("login_test_current_user.txt")
	defer f.Close()
	defer func() {
		terminalReadPassword = terminal.ReadPassword
	}()
	scanner := bufio.NewScanner(f)
	terminalReadPassword = func (fd int) ([]byte, error) {
		scanner.Scan()
		return scanner.Bytes(), nil
	}
	assert.Equal(t, Login(scanner, ""), "bob")
}
