package server

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"chat/protocol"
	"chat/db"
)

func getFakeUser() *db.User {
	u := new(db.User)
	u.IP = "1.2.3.4"
	u.MAC = "ab:cd:ef:gh:ij"
	u.Username = "sam"
	return u
}

const (
	message = "Hello World"
)

func TestNewMessage(t *testing.T) {
	u := getFakeUser()
	m := NewMessage(u, u.IP, message)
	assert.Equal(t, u.IP, m.SourceIP)
	assert.Equal(t, u.MAC, m.SourceMAC)
	assert.Equal(t, message, m.Text)
}

func TestMessage_StartProtocol(t *testing.T) {
	u := getFakeUser()
	m := NewMessage(u, u.IP, message)
	proto := protocol.OTRProtocol{}
	m.StartProtocol(proto)

	assert.Equal(t, proto.ToType(), m.StartProto)
}
