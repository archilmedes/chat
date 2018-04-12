package db

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func ConversationTest (t *testing.T){
	conversations := QueryConversations()
	assert.Equal(t, 8, len(conversations))
	assert.Equal(t, 52, conversations[3].SSID)
	assert.Equal(t, 0, conversations[3].sentOrReceived)
	assert.Equal(t, "03/28/2018:18:04:10", conversations[3].timestamp)
	assert.Equal(t, "lul", conversations[3].message)
}

