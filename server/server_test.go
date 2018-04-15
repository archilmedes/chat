package server

import (
	"chat/core"
	"chat/protocol"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func startUpServer(t *testing.T) Server {
	var server Server
	mac, ip, _ := core.GetAddresses()
	assert.NoError(t, server.Start("Archil", mac, ip))
	return server
}

func setUpServerAndHandshake(t *testing.T) Server {
	server := startUpServer(t)
	assert.NotNil(t, server.User)
	// Initialize a session with yourself
	err := server.StartSession(server.User.IP, protocol.NewOTRProtocol())
	assert.Nil(t, err)
	assert.Nil(t, err)
	// Let time pass for handshake to complete
	time.Sleep(2000 * time.Millisecond)
	return server
}

func TestServer_Start(t *testing.T) {
	server := startUpServer(t)
	assert.NotEqual(t, nil, server.Listener)
	assert.NoError(t, server.Shutdown())
}

func TestServer_Send(t *testing.T) {
	server := setUpServerAndHandshake(t)
	assert.NoError(t, server.Send(server.User.IP, "Hello World!"))
	time.Sleep(1 * time.Second)
	server.Shutdown()
}

func TestServer_GetSessionsToIP(t *testing.T) {
	server := setUpServerAndHandshake(t)
	sessions := server.GetSessionsToIP(server.User.IP)
	assert.Equal(t, 2, len(sessions))

	msg := []byte("Hello world")
	cyp, _ := sessions[0].Proto.Encrypt(msg)

	msgBack, _ := sessions[1].Proto.Decrypt(cyp[0])
	assert.Equal(t, msgBack[0], msg)
}
