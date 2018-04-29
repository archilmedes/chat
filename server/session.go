package server

import (
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/protocol"
	"time"
	"github.com/wavyllama/chat/core"
)

// Struct for a messaging session between a user and his/her friend
type Session struct {
	From      *db.User
	To        *db.Friend
	Proto     protocol.Protocol
	StartTime time.Time
}

// Return a new session between a user and his/her friend with a protocol
func NewSession(from *db.User, to *db.Friend, protocol protocol.Protocol) *Session {
	session := new(Session)
	(*session).From = from
	(*session).To = to
	(*session).Proto = protocol
	(*session).StartTime = time.Now()
	return session
}

// Return a new session between a user and his/her friend based on a message
func NewSessionFromUserAndMessage(from *db.User, to *db.Friend, protoType string) *Session {
	return NewSession(from, to, protocol.CreateProtocolFromType(protoType))
}

// Ends the current session
func (s *Session) EndSession() bool {
	s.Proto.EndSession()
	return db.DeleteSession(s.Proto.GetSessionID())
}

// Saves a session to the database
func (s *Session) Save() bool {
	sessionID := s.Proto.GetSessionID()
	bb := s.Proto.Serialize()

	return db.InsertIntoSessions(sessionID, s.From.Username, s.To.DisplayName, s.Proto.ToType(), bb, core.GetFormattedTime(s.StartTime))
}
