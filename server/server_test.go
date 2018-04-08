package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const Name = "Generic"

func TestServer_Start(t *testing.T) {
	var program Server
	assert.NoError(t, program.Start(Name))
	assert.Equal(t, Name, program.Username)
	assert.NotEqual(t, nil, program.Listener)
	assert.NotEqual(t, "", program.IP)
	assert.NoError(t, program.Shutdown())
}

func TestServer_Send(t *testing.T) {
	var program Server
	program.Start(Name)
	assert.NoError(t, program.Send(program.IP, "Hello World!"))
	program.Shutdown()
}
