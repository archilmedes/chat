package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/protocol"
	"testing"
)

const (
	bobDisplayName = "Bobby B"
)

func setUpSession() (*db.User, *db.Friend, *Session) {
	db.SetupEmptyTestDatabase()

	alice := new(db.User)
	alice.IP = "1.2.3.4"
	bob := new(db.Friend)
	bob.IP = "5.6.7.8"
	bob.DisplayName = bobDisplayName
	sess := NewSession(alice, bob, protocol.NewOTRProtocol())
	return alice, bob, sess
}

func TestSession_EndSession(t *testing.T) {
	alice, _, sess := setUpSession()
	assert.True(t, sess.Save())
	sessions := alice.GetSessions(bobDisplayName)
	assert.Equal(t, 1, len(sessions))
	assert.Equal(t, sess.Proto.GetSessionID(), sessions[0].SSID)
	assert.Equal(t, sess.Proto.Serialize(), sessions[0].ProtocolValue)

	sess.EndSession()
	assert.False(t, sess.Proto.IsActive())
}
