package core

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestRejectFriend(t *testing.T) {
	os.Stdin.Close()
	f, _ := os.Open("input_test_reject.txt")
	CondWait = func() {
		Friending = REJECT
	}
	defer f.Close()
	CondWait = func() {
		Friending = ACCEPT
	}
	bufioNewScanner = func(r io.Reader) *bufio.Scanner {
		return bufio.NewScanner(f)
	}
	defer func() {
		CondWait = Cond.Wait
		bufioNewScanner = bufio.NewScanner
	}()
	assert.Equal(t, "", GetDisplayNameFromConsole("", ""))
}

func TestAcceptFriend(t *testing.T) {
	os.Stdin.Close()
	f, _ := os.Open("input_test_accept.txt")
	defer f.Close()
	CondWait = func() {
		Friending = ACCEPT
	}
	bufioNewScanner = func(r io.Reader) *bufio.Scanner {
		return bufio.NewScanner(f)
	}
	defer func() {
		CondWait = Cond.Wait
		bufioNewScanner = bufio.NewScanner
	}()
	assert.Equal(t, "archie", GetDisplayNameFromConsole("", ""))
}
