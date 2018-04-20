package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func ConversationsTest(t *testing.T) {
	ConversationsTestID(t)
	ConversationsTestSSID(t)
}

func ConversationsTestID(t *testing.T) {
	conversations := GetConversationUsers("karateAMD", "10.192.345.987")
	assert.Equal(t, 3, len(conversations))
	assert.Equal(t, 34, int(conversations[0].Session.SSID))
	assert.Equal(t, 34, int(conversations[1].Message.SSID))
	assert.Equal(t, 34, int(conversations[2].Session.SSID))
	assert.Equal(t, "2018-04-10 12:30:08.222222", conversations[0].Message.Timestamp)
	assert.Equal(t, "When are we playing Fortnite?", conversations[2].Message.Text)
}

func ConversationsTestSSID(t *testing.T) {
	conversations := GetConversationSSID(34)
	assert.Equal(t, 3, len(conversations))
	assert.Equal(t, 34, int(conversations[0].Session.SSID))
	assert.Equal(t, "plain", conversations[0].Session.ProtocolType)
	assert.Equal(t, "I almost made my Mac a brick", conversations[1].Message.Text)
	assert.Equal(t, 34, int(conversations[2].Session.SSID))
	assert.Equal(t, "karateAMD", conversations[2].Session.Username)
	assert.Equal(t, "10.192.345.987", conversations[2].Session.FriendMac)
	assert.Equal(t, "2018-04-08 17:59:02.777777", conversations[2].Message.Timestamp)
	assert.Equal(t, "plain", conversations[2].Session.ProtocolType)
}
