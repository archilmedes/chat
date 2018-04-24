package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSessions(t *testing.T) {
	SetupDatabaseForTests(t)
	InitialSetupTest(t)
	InsertSessionTest(t)
	DeleteSessionTest(t)
}

func InitialSetupTest(t *testing.T) {
	sessions := QuerySessions()
	assert.Equal(t, 7, len(sessions))
	assert.Equal(t, "plain", sessions[1].ProtocolType)
	assert.Equal(t, []byte("lastLine7"), sessions[2].ProtocolValue)
	assert.Equal(t, "karateAMD", sessions[3].Username)
	assert.Equal(t, uint64(35), sessions[4].SSID)
	assert.Equal(t, "123.456.789", sessions[5].FriendDisplayName)
	assert.Equal(t, "otr", sessions[6].ProtocolType)
	assert.Equal(t, []byte("number5"), sessions[6].ProtocolValue)
	assert.Equal(t, "andrew", sessions[6].Username)
	assert.Equal(t, uint64(64), sessions[6].SSID)
	assert.Equal(t, "10.192.345.987", sessions[6].FriendDisplayName)
}

func InsertSessionTest(t *testing.T) {
	assert.True(t, InsertIntoSessions(84, "bill", "123.333.333.456", "otr", []byte("newStringVal"), "2018-02-02 02:03:04.567890"))
	sessions := QuerySessions()
	assert.Equal(t, 8, len(sessions))
	assert.Equal(t, uint64(84), sessions[7].SSID)
	assert.Equal(t, "123.333.333.456", sessions[7].FriendDisplayName)
	assert.Equal(t, []byte("newStringVal"), sessions[7].ProtocolValue)
}

func DeleteSessionTest(t *testing.T) {
	assert.True(t, DeleteSession(81))
	sessions := QuerySessions()
	assert.Equal(t, 8, len(sessions))
}
