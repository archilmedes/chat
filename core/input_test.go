package core

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	"time"
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

func TestGetFormattedTime(t *testing.T) {
	someTime := time.Unix(100, 1000)
	assert.Equal(t, "1969-12-31 18:01:40.000001", GetFormattedTime(someTime))
}
