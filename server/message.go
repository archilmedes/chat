package server

import "chat/protocol"

// Struct for messages being sent
type Message struct {
	MAC, IP string
	StartProto protocol.Protocol
	Text []byte
}

// Create new message to send a message
func NewMessage(mac string, ipAddress string, text []byte) (*Message) {
	m := Message{MAC: mac, IP: ipAddress, Text: text}
	return &m
}

func (m *Message) StartProtocol(proto protocol.Protocol) {
	m.StartProto = proto
}