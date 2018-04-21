package server

import (
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/protocol"
	"time"
	"log"
)

// Struct for a messaging session between a user and his/her friend
type Session struct {
	From      *db.User
	To        *db.Friend
	Proto     protocol.Protocol
	StartTime time.Time
}

// Return a new session between a user and his/her friend with a protocol
func NewSession(from *db.User, to *db.Friend, protocol protocol.Protocol, startTime time.Time) *Session {
	session := new(Session)
	(*session).From = from
	(*session).To = to
	(*session).Proto = protocol
	(*session).StartTime = startTime
	return session
}

// Return a new session between a user and his/her friend based on a message
func NewSessionFromUserAndMessage(from *db.User, to *db.Friend, protoType string, startSessionTime time.Time) *Session {
	return NewSession(from, to, protocol.CreateProtocolFromType(protoType), startSessionTime)
}

// Ends the current session
func (s *Session) EndSession() bool {
	s.Proto.EndSession()
	return db.DeleteSession(s.Proto.GetSessionID())
}

// Return all messages that have been sent between two users in a given this session
func (s *Session) GetMessages() [][]byte {
	converse := db.GetConversationUsers(s.From.Username, s.To.DisplayName)
	var messages [][]byte
	for _, c := range converse {
		dec, err := s.Proto.Decrypt([]byte(c.Message.Text))
		if err != nil {
			log.Printf("GetMessages: %s", err.Error())
		}
		messages = append(messages, dec[0])
	}
	return messages
}

// Saves a session to the database
func (s *Session) Save() {
	sessionID := s.Proto.GetSessionID()
	db.InsertIntoSessions(sessionID, s.From.Username, s.To.MAC, s.Proto.ToType(), string(s.Proto.Serialize()), s.StartTime.Format(time.RFC3339))
}
