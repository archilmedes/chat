package server

import "chat/protocol"

// Struct for messages being sent
type Message struct {
	SourceMAC, SourceIP, DestIP string
	StartProto                  protocol.Protocol // If a protocol is started, this will be defined
	Text                        string
}

// Create new message to send a message
func NewMessage(from *User, destIp string, text string) (*Message) {
	m := Message{
		SourceMAC: (*from).MAC, SourceIP: (*from).IP,
		DestIP: destIp, Text: text,
		StartProto: protocol.PlainProtocol{}}
	return &m
}

func (m *Message) StartProtocol(proto protocol.Protocol) {
	m.StartProto = proto
}