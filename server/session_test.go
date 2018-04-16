package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/protocol"
	"testing"
	"time"
)

func TestSession_ConverseWith(t *testing.T) {
	alice, bob, sess := setUpSession()

	assert.True(t, sess.ConverseWith(bob.IP))
	assert.False(t, sess.ConverseWith(alice.IP))
}

func setUpSession() (*db.User, *db.Friend, *Session) {
	alice := new(db.User)
	alice.IP = "1.2.3.4"
	bob := new(db.Friend)
	bob.IP = "5.6.7.8"
	sess := NewSession(alice, bob, protocol.NewOTRProtocol(), time.Now())
	return alice, bob, sess
}

func TestSession_EndSession(t *testing.T) {
	_, _, sess := setUpSession()
	sess.EndSession()
	assert.False(t, sess.Proto.IsActive())
}
