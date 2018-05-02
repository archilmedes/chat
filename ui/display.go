package ui

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
	Sender, Time string
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

func (m *InfoMessage) New(info string) {
	m.createInfoMessage(info)
}

func (m *ReceiveChat) New(info, sender, time string) {
	m.createInfoMessage(info)
	m.Sender = sender
	m.Time = time
}

func (m *FriendRequest) New(info, username, ip string) {
	m.createInfoMessage(info)
	m.Username = username
	m.IP = ip
}
