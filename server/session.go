package server

import (
	"chat/db"
	"chat/protocol"
	"time"
)

type Session struct {
	From      *db.User
	To        *db.Friend
	Proto     protocol.Protocol
	StartTime time.Time
}

// Return a new session between a user and their friend with a protocol
func NewSession(from *db.User, to *db.Friend, protocol protocol.Protocol, startTime time.Time) *Session {
	session := new(Session)
	(*session).From = from
	(*session).To = to
	(*session).Proto = protocol
	(*session).StartTime = startTime
	return session
}

func NewSessionFromUserAndMessage(from *db.User, msg Message) *Session {
	friend := new(db.Friend)
	friend.IP = msg.SourceIP
	friend.MAC = msg.SourceMAC
	return NewSession(from, friend, protocol.CreateProtocolFromType(msg.StartProto), msg.StartProtoTimestamp)
}

// Ends the current session
func (s *Session) EndSession() {
	s.Proto.EndSession()
}

// Returns true if the session is conversing with a use defined by their SourceIP address
func (s *Session) ConverseWith(sourceIp string) bool {
	return (*s.To).IP == sourceIp
}
