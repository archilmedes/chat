package server

import "chat/protocol"

// Struct for messages being sent
type Message struct {
	SourceMAC, SourceIP, DestIP string
	StartProto                  protocol.Protocol // If a protocol is started, this will be defined
	Text                        []byte
}

// Create new message to send a message
func NewMessage(from *User, destIp string, text []byte) (*Message) {
	m := Message{SourceMAC: (*from).MAC, SourceIP: (*from).IP, DestIP: destIp, Text: text}
	return &m
}

func (m *Message) StartProtocol(proto protocol.Protocol) {
	m.StartProto = proto
}