package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"chat/core"
	"chat/protocol"
)

func startUpServer(t *testing.T) Server {
	var server Server
	mac, ip, _ := core.GetAddresses()
	assert.NoError(t, server.Start("Archil", mac, ip))
	return server
}

func TestServer_Start(t *testing.T) {
	server := startUpServer(t)
	assert.NotEqual(t, nil, server.Listener)
	assert.NoError(t, server.Shutdown())
}

func TestServer_Send(t *testing.T) {
	server := startUpServer(t)
	assert.NotNil(t, server.User)
	// Send a message to yourself
	server.StartSession(server.User.IP, protocol.NewOTRProtocol())

	assert.NoError(t, server.Send(server.User.IP, "Hello World!"))
	server.Shutdown()
}

func TestServer_CreateOrGetSession_create(t *testing.T) {
	server := startUpServer(t)

	u := getFakeUser()
	msg := []byte("Hello world")
	message := NewMessage(u, u.IP, string(msg))
	message.StartProto = protocol.NewOTRProtocol()

	sess := server.CreateOrGetSession(*message)
	assert.NotNil(t, sess)
	assert.Equal(t, server.User, sess.From)
	assert.Equal(t, u.IP, sess.To.IP)
	assert.Equal(t, message.StartProto, sess.Proto)
	server.Shutdown()
}

func TestServer_CreateOrGetSession_get(t *testing.T) {
	server := startUpServer(t)

	u := getFakeUser()
	msg := []byte("Hello world")
	message := NewMessage(u, u.IP, string(msg))
	message.StartProto = protocol.OTRProtocol{}

	msg2 := []byte("Hello you")
	message2 := NewMessage(u, u.IP, string(msg2))

	sess := server.CreateOrGetSession(*message)
	assert.NotNil(t, sess)

	sess2 := server.CreateOrGetSession(*message2)
	assert.NotNil(t, sess2)

	assert.Equal(t, sess, sess2)
	server.Shutdown()
}
