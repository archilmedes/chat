package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChatMessage_Init(t *testing.T) {
	var msg Message
	mac := "12:34:56:78:90"
	text := []byte("Hello World!")
	msg.Init(mac, text)
	assert.Equal(t, mac, msg.MAC)
	assert.Equal(t, text, msg.Text)
}
