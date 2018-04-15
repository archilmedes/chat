package server

import (
	"chat/protocol"
	"chat/db"
)

type Session struct {
	From, To *db.User
	Proto *protocol.Protocol
}

// Return a new session between two users with a protocol
func NewSession(from *db.User, to *db.User, protocol *protocol.Protocol) (*Session) {
	session := new(Session)
	(*session).From = from
	(*session).To = to
	(*session).Proto = protocol
	return session
}

// Ends the current session
func (s *Session) EndSession() {
	(*s.Proto).EndSession()
}