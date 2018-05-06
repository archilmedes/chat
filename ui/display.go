package ui

import "time"

// Any message meant to be displayed in UI
type DisplayMessage interface {
	Body() string
}

// Message of just text information
type InfoMessage struct {
	DisplayMessage
	Message string
}

// Chat message
type ReceiveChat struct {
	InfoMessage
	Sender string
	Time time.Time
}

// Friend request received
type FriendRequest struct {
	InfoMessage
	Username, IP string
}

func (m *InfoMessage) Body() string {
	return m.Message
}

func (m *InfoMessage) createInfoMessage(info string) {
	m.Message = info
}

func NewInfoMessage(info string) *InfoMessage {
	return &InfoMessage{Message: info}
}

func NewReceiveChatMessage(info, sender string, time time.Time) *ReceiveChat {
	return &ReceiveChat{InfoMessage: *NewInfoMessage(info),
						Sender: sender, Time: time}
}

func NewFriendRequestMessage(info, username, ip string) *FriendRequest {
	return &FriendRequest{InfoMessage: *NewInfoMessage(info),
							Username: username, IP: ip}
}