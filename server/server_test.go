package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"chat/core"
)

func TestServer_Start(t *testing.T) {
	var program Server
	mac, ip, _ := core.GetAddresses()
	assert.NoError(t, program.Start("Archil", mac, ip))
	assert.NotEqual(t, nil, program.Listener)
	assert.NoError(t, program.Shutdown())
}

func TestServer_Send(t *testing.T) {
	var program Server
	mac, ip, _ := core.GetAddresses()
	assert.NoError(t, program.Start("Archil", mac, ip))
	assert.NotNil(t, program.User)
	assert.NoError(t, program.Send(program.User.IP, []byte("Hello World!")))
	program.Shutdown()
}
