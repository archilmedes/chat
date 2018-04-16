package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func MessagesTest(t *testing.T) {
	MessagesSetup(t)
	GetConversationTest(t)
}
func GetConversationTest(t *testing.T) {
	conversations := GetConversation(3, 5)
	assert.Equal(t, 3, len(conversations))
	assert.Equal(t, 35, conversations[0].SSID)
	assert.Equal(t, 35, conversations[1].SSID)
	assert.Equal(t, 35, conversations[2].SSID)
	assert.Equal(t, "04/10/2018:12:30:08", conversations[0].timestamp)
	assert.Equal(t, "When are we playing Fortnite?", conversations[2].message)
}

func MessagesSetup(t *testing.T) {
	messages := QueryMessages()
	assert.Equal(t, 8, len(messages))
	assert.Equal(t, 52, messages[3].SSID)
	assert.Equal(t, 0, messages[3].sentOrReceived)
	assert.Equal(t, "03/28/2018:18:04:10", messages[3].timestamp)
	assert.Equal(t, "lul", messages[3].message)
}
