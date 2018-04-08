package server

// Struct for messages being sent
type Message struct {
	User, Text string
}

// Create message
func (m *Message) Init(user string, text string) {
	(*m).User = user
	(*m).Text = text
}