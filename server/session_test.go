package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/protocol"
	"testing"
	"time"
)

func setUpSession() (*db.User, *db.Friend, *Session) {
	alice := new(db.User)
	alice.IP = "1.2.3.4"
	bob := new(db.Friend)
	bob.IP = "5.6.7.8"
	sess := NewSession(alice, bob, protocol.NewOTRProtocol(), time.Now())

	dbDeleteSession = func(SSID uint64) bool {
		return true
	}
	return alice, bob, sess
}

func TestSession_EndSession(t *testing.T) {
	_, _, sess := setUpSession()
	defer func() {
		dbDeleteSession = db.DeleteSession
	}()
	sess.EndSession()
	assert.False(t, sess.Proto.IsActive())
}
