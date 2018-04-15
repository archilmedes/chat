// Package server provides a server implementation to send messages
// between two different servers running on different instances of the app
package server

import (
	"chat/db"
	"chat/protocol"
	"time"
)

// Struct for messages being sent
type Message struct {
	SourceMAC, SourceIP, DestIP string
	StartProtoTimestamp         time.Time
	StartProto, Text            string // If a protocol is started, StartProto will be defined
	ID                          int
	Handshake                   bool
}

// Create new message to send a message
func NewMessage(from *db.User, destIp string, text string) *Message {
	m := Message{
		SourceMAC: (*from).MAC, SourceIP: (*from).IP,
		DestIP: destIp, Text: text,
		StartProto: ""}
	return &m
}

// Start a protocol and its handshake
func (m *Message) StartProtocol(proto protocol.Protocol) {
	m.StartProtoTimestamp = time.Now()
	m.Handshake = true
	m.StartProto = proto.ToType()
}
