package server

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"chat/protocol"
)

func getFakeUser() *User {
	u := new(User)
	u.IP = "1.2.3.4"
	u.MAC = "ab:cd:ef:gh:ij"
	u.Username = "sam"
	return u
}

func TestNewMessage(t *testing.T) {
	u := getFakeUser()
	msg := string([]byte("Hello world"))
	m := NewMessage(u, u.IP, msg)
	assert.Equal(t, u.IP, m.SourceIP)
	assert.Equal(t, u.MAC, m.SourceMAC)
	assert.Equal(t, msg, m.Text)
}

func TestMessage_StartProtocol(t *testing.T) {
	u := getFakeUser()
	msg := string([]byte("Hello world"))
	m := NewMessage(u, u.IP, msg)
	proto := protocol.OTRProtocol{}
	m.StartProtocol(proto)

	assert.Equal(t, proto, m.StartProto)
}
