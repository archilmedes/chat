// Package server provides a server implementation to send messages
// between two different servers running on different instances of the app
package server

import (
	"time"
)

// Generic interface for messages being sent and received
type Message interface {
	SourceID() (string, string) // MAC address, Username
	DestID() string             // Username
}

type GenericMessage struct {
	Message
	SourceMAC, SourceUsername, DestUsername string
}

// Message for sending and receiving friend requests/info
type FriendMessage struct {
	GenericMessage
}

// Message for handshaking an securing a session
type HandshakeMessage struct {
	GenericMessage
	Round       int
	SessionTime time.Time
	ProtoType   string
	Secret      []byte
}

// Message for sending regular information to a friend
type ChatMessage struct {
	GenericMessage
	Text []byte
}

func (m *GenericMessage) NewPayload(SourceMAC, SourceUsername, DestUsername string) {
	(*m).SourceMAC = SourceMAC
	(*m).SourceUsername = SourceUsername
	(*m).DestUsername = DestUsername
}

// Get source-identifying MAC and username info
func (m *GenericMessage) SourceID() (string, string) {
	return (*m).SourceMAC, (*m).SourceUsername
}

// Get destination-identifying username
func (m *GenericMessage) DestID() string {
	return (*m).DestUsername
}
