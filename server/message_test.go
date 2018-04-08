package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChatMessage_Init(t *testing.T) {
	var msg Message
	user := "Archil"
	text := "Hello World!"
	msg.Init(user, text)
	assert.Equal(t, user, msg.User)
	assert.Equal(t, text, msg.Text)
}
