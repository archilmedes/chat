package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func MessagesTest(t *testing.T) {
	MessagesSetup(t)
}

func MessagesSetup(t *testing.T) {
	messages := QueryMessages()
	assert.Equal(t, 8, len(messages))
	assert.Equal(t, 52, messages[3].SSID)
	assert.Equal(t, 0, messages[3].sentOrReceived)
	assert.Equal(t, "03/28/2018:18:04:10", messages[3].timestamp)
	assert.Equal(t, "lul", messages[3].message)
}
