package db

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func ConversationsTest(t *testing.T) {
	ConversationsTestID(t)
	ConversationsTestSSID(t)
}

func ConversationsTestID(t *testing.T) {
	conversations := GetConversationID(3, 5)
	assert.Equal(t, 3, len(conversations))
	assert.Equal(t, 35, conversations[0].Session.SSID)
	assert.Equal(t, 35, conversations[1].Message.SSID)
	assert.Equal(t, 35, conversations[2].Session.SSID)
	assert.Equal(t, "04/10/2018:12:30:08", conversations[0].Message.timestamp)
	assert.Equal(t, "When are we playing Fortnite?", conversations[2].Message.message)
}

func ConversationsTestSSID(t *testing.T) {
	conversations := GetConversationSSID(35)
	assert.Equal(t, 3, len(conversations))
	assert.Equal(t, 35, conversations[0].Session.SSID)
	assert.Equal(t, "abcdb378675934bdbd0935847349036985fbd490590584374374b43894784578431243b37465723894d3434981fdcb484726739923874bd3837473fedb31", conversations[0].Session.PrivateKey)
	assert.Equal(t, "I almost made my Mac a brick", conversations[1].Message.message)
	assert.Equal(t, 35, conversations[2].Session.SSID)
	assert.Equal(t, 3, conversations[2].Session.UserId)
	assert.Equal(t, 5, conversations[2].Session.FriendId)
	assert.Equal(t, "04/08/2018:17:59:02", conversations[2].Message.timestamp)
	assert.Equal(t, "ACD537492B126CD43", conversations[2].Session.Fingerprint)
}
