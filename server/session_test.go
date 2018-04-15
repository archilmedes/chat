package server

import (
	"chat/db"
	"chat/protocol"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSession_ConverseWith(t *testing.T) {
	alice, bob, sess := setUpSession()

	assert.True(t, sess.ConverseWith(bob.IP))
	assert.False(t, sess.ConverseWith(alice.IP))
}

func setUpSession() (*db.User, *Friend, *Session) {
	alice := new(db.User)
	alice.IP = "1.2.3.4"
	bob := new(Friend)
	bob.IP = "5.6.7.8"
	sess := NewSession(alice, bob, protocol.NewOTRProtocol(), time.Now())
	return alice, bob, sess
}

func TestSession_EndSession(t *testing.T) {
	_, _, sess := setUpSession()
	sess.EndSession()
	assert.False(t, sess.Proto.IsActive())
}
