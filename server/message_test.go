package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/protocol"
	"testing"
)

const MessageText = "Hello World"

func getFakeUser() *db.User {
	u := new(db.User)
	u.IP = "1.2.3.4"
	u.MAC = "ab:cd:ef:gh:ij"
	u.Username = "sam"
	return u
}

func TestNewMessage(t *testing.T) {
	u := getFakeUser()
	m := NewMessage(u, u.IP, MessageText)
	assert.Equal(t, u.IP, m.SourceIP)
	assert.Equal(t, u.MAC, m.SourceMAC)
	assert.Equal(t, MessageText, m.Text)
}

func TestMessage_StartProtocol(t *testing.T) {
	u := getFakeUser()
	m := NewMessage(u, u.IP, MessageText)
	proto := protocol.OTRProtocol{}
	m.StartProtocol(proto)

	assert.Equal(t, proto.ToType(), m.StartProto)
}
