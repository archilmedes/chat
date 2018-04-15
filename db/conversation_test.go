package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func ConversationTest(t *testing.T) {
	ConversationSetup(t)
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

func ConversationSetup(t *testing.T) {
	conversations := QueryConversations()
	assert.Equal(t, 8, len(conversations))
	assert.Equal(t, 52, conversations[3].SSID)
	assert.Equal(t, 0, conversations[3].sentOrReceived)
	assert.Equal(t, "03/28/2018:18:04:10", conversations[3].timestamp)
	assert.Equal(t, "lul", conversations[3].message)
}
