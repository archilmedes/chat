package server

import (
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/protocol"
	"time"
)

const(
	bobDisplayName = "Bobby B"
)

func setUpSession() (*db.User, *db.Friend, *Session) {
	db.SetupEmptyTestDatabase()

	alice := new(db.User)
	alice.IP = "1.2.3.4"
	bob := new(db.Friend)
	bob.IP = "5.6.7.8"
	bob.DisplayName = bobDisplayName
	sess := NewSession(alice, bob, protocol.NewOTRProtocol(), time.Now())
	return alice, bob, sess
}