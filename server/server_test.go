package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServer_Start(t *testing.T) {
	var program Server
	assert.NoError(t, program.Start("Archil", "random", "192.168.86.22"))
	assert.NotEqual(t, nil, program.Listener)
	assert.NoError(t, program.Shutdown())
}

func TestServer_Send(t *testing.T) {
	var program Server
	program.User.IP = "1.2.3.4"
	assert.NoError(t, program.Send(program.User.IP, "12:34:56:78:90", []byte("Hello World!")))
	program.Shutdown()
}
