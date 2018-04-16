package core

import (
	"bufio"
	"chat/db"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	// TODO: More test with delete user
	db.SetupDatabase()
	os.Stdin.Close()
	f, _ := os.Open("login_test.txt")
	defer f.Close()
	defer func() {
		terminalReadPassword = terminal.ReadPassword
	}()
	scanner := bufio.NewScanner(f)
	terminalReadPassword = func(fd int) ([]byte, error) {
		scanner.Scan()
		return scanner.Bytes(), nil
	}
	assert.Equal(t, Login(scanner, ""), "bob")
}
