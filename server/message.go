package server

// Struct for messages being sent
type Message struct {
	MAC string
	Text []byte
}

// Create message
func (m *Message) Init(mac string, text []byte) {
	(*m).MAC = mac
	(*m).Text = text
}