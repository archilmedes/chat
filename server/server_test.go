package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/wavyllama/chat/core"
	"testing"
	"github.com/wavyllama/chat/db"
	"time"
)

func startUpServer(t *testing.T) Server {
	var server Server
	mac, ip, _ := core.GetAddresses()
	db.SetupEmptyTestDatabase()
	assert.NoError(t, server.Start("Archil", mac, ip))
	// Let time pass for handshake to complete
	time.Sleep(2000 * time.Millisecond)
	return server
}

func TestServer_Start(t *testing.T) {
	server := startUpServer(t)
	assert.NotEqual(t, nil, server.Listener)
	assert.NoError(t, server.Shutdown())
}

func TestServer_GetSessionsWithFriend(t *testing.T) {
	server := startUpServer(t)
	sessions := server.GetSessionsWithFriend(server.User.MAC, server.User.Username)
	assert.Equal(t, 2, len(sessions))

	msg := []byte("Hello world")
	cyp, _ := sessions[0].Proto.Encrypt(msg)

	msgBack, _ := sessions[1].Proto.Decrypt(cyp[0])
	assert.Equal(t, msgBack[0], msg)
}
