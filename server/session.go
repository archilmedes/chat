package server

import (
	"chat/protocol"
)

type Session struct {
	From, To *User
	Proto protocol.Protocol
}

// Return a new session between two users with a protocol
func NewSession(from *User, to *User, protocol protocol.Protocol) (*Session) {
	session := new(Session)
	(*session).From = from
	(*session).To = to
	(*session).Proto = protocol
	return session
}

// Ends the current session
func (s *Session) EndSession() {
	s.Proto.EndSession()
}

// Returns true if the session is conversing with a use defined by their IP address
func (s *Session) ConverseWith(destIp string) bool {
	return (*s.To).IP == destIp
}