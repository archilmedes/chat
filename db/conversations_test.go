package db

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func ConversationsTest(t *testing.T) {
	conversations := GetConversation(3, 5)
	assert.Equal(t, 3, len(conversations))
	assert.Equal(t, 35, conversations[0].Session.SSID)
	assert.Equal(t, 35, conversations[1].Message.SSID)
	assert.Equal(t, 35, conversations[2].Session.SSID)
	assert.Equal(t, "04/10/2018:12:30:08", conversations[0].Message.timestamp)
	assert.Equal(t, "When are we playing Fortnite?", conversations[2].Message.message)
}
