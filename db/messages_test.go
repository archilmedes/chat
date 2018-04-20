package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func MessagesTest(t *testing.T) {
	MessagesSetup(t)
	InsertMessageTest(t)
	DeleteMessageTest(t)
}

func MessagesSetup(t *testing.T) {
	messages := QueryMessages()
	assert.Equal(t, 8, len(messages))
	assert.Equal(t, 52, messages[3].SSID)
	assert.Equal(t, 0, messages[3].SentOrReceived)
	assert.Equal(t, "2018-03-28 18:04:10.333333", messages[3].Timestamp)
	assert.Equal(t, "lul", messages[3].Text)
}

func InsertMessageTest(t *testing.T) {
	assert.True(t, InsertMessage(52, "wassup", "2018-04-12 05:01:10.888888", Received))
	messages := QueryMessages()
	assert.Equal(t, 9, len(messages))
	assert.Equal(t, "wassup", messages[8].Text)
}

func DeleteMessageTest(t *testing.T) {
	assert.True(t, DeleteMessage(52, "wassup", "2018-04-12 05:01:10.888888", Received))
	messages := QueryMessages()
	assert.Equal(t, 8, len(messages))
}
