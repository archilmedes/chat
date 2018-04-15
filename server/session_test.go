package server

import (
	"testing"
	"chat/protocol"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestSession_ConverseWith(t *testing.T) {
	_, bob, sess := setUpSession()

	assert.True(t, sess.ConverseWith(bob.IP))
}

func setUpSession() (*User, *Friend, *Session) {
	alice := new(User)
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

func TestNewSession(t *testing.T) {
	alice := new(User)
	alice.IP = "1.2.3.4"

	a2 := new(Friend)
	a2.IP = alice.IP

	t1 := time.Now()
	// 1. alice -> a2 at time t1
	// a2 creates session between a2 and alice with t1 time b/c no session exists
	NewSession(alice, a2, protocol.NewOTRProtocol(), t1)

	// 2. a2 -> alice at time t1 + 10ms
	// alice creates session between alice and a2 with t1 + 10ms time b/c no session exists with t1 + 10ms
	NewSession(alice, a2, protocol.NewOTRProtocol(), t1.Add(10 * time.Millisecond))


	// TODO check from -> to and to -> from, if there is < 2, create one, otherwise just return he
	// 3. alice -> a2 at time t1 + 20ms
	// a2 finds existing session of time

	// alice gets it, checks for session, and if it sees one that doesn't have same timestamp, then that's the right one
}