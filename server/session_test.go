package server

import (
	"testing"
	"chat/protocol"
	"github.com/stretchr/testify/assert"
)

func TestSession_ConverseWith(t *testing.T) {
	_, bob, sess := setUpSession()

	assert.True(t, sess.ConverseWith(bob.IP))
}

func setUpSession() (*User, *User, *Session) {
	alice := new(User)
	alice.IP = "1.2.3.4"
	bob := new(User)
	bob.IP = "5.6.7.8"
	sess := NewSession(alice, bob, protocol.NewOTRProtocol())
	return alice, bob, sess
}

func TestSession_EndSession(t *testing.T) {
	_, _, sess := setUpSession()
	sess.EndSession()
	assert.False(t, sess.Proto.IsActive())
}