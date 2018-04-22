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

// Message for sending and receiving friend requests/info
type FriendMessage struct {
	SourceMAC, SourceIPAddress, SourceUsername, DestUsername string
}

// Message for handshaking an securing a session
type HandshakeMessage struct {
	Round                                                      int
	SessionTime                                                time.Time
	SourceMAC, SourceUsername, DestUsername, Secret, Prototype string
}

// Message for sending regular information to a friend
type ChatMessage struct {
	SourceMAC, SourceUsername, DestUsername, Text string
}

// Same implementations for getting ID for sender and receiver:

func (m *FriendMessage) SourceID() (string, string) {
	return (*m).SourceMAC, (*m).SourceUsername
}

func (m *HandshakeMessage) SourceID() (string, string) {
	return (*m).SourceMAC, (*m).SourceUsername
}

func (m *ChatMessage) SourceID() (string, string) {
	return (*m).SourceMAC, (*m).SourceUsername
}

func (m *FriendMessage) DestID() string {
	return (*m).DestUsername
}

func (m *HandshakeMessage) DestID() string {
	return (*m).DestUsername
}

func (m *ChatMessage) DestID() string {
	return (*m).DestUsername
}
