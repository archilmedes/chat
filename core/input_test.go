package core

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func setupTests() *os.File {
	os.Stdin.Close()
	f, _ := os.Open("input_test.txt")
	bufioNewScanner = func(r io.Reader) *bufio.Scanner {
		return bufio.NewScanner(f)
	}
	return f
}

func TestRejectFriend(t *testing.T) {
	f := setupTests()
	defer f.Close()
	CondWait = func() {
		Friending = REJECT
	}
	defer func() {
		CondWait = Cond.Wait
		bufioNewScanner = bufio.NewScanner
	}()
	assert.Equal(t, "", GetDisplayNameFromConsole("", ""))
}

func TestAcceptFriend(t *testing.T) {
	f := setupTests()
	defer f.Close()
	CondWait = func() {
		Friending = ACCEPT
	}
	defer func() {
		CondWait = Cond.Wait
		bufioNewScanner = bufio.NewScanner
	}()
	assert.Equal(t, "archie", GetDisplayNameFromConsole("", ""))
}
