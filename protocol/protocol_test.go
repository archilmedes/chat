package protocol

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var message = "testMessage"

func TestPlainProtocol_Encrypt(t *testing.T) {
	p := new(PlainProtocol)
	mes := []byte(message)
	res, err := p.Encrypt(mes)
	assert.Nil(t, err)
	assert.Equal(t, res[0], mes)
}

func TestPlainProtocol_Decrypt(t *testing.T) {
	p := new(PlainProtocol)
	cypher := []byte(message)
	msg, err := p.Decrypt(cypher)
	assert.Nil(t, err)
	assert.Equal(t, msg[0], cypher)
}

func TestPlainProtocol_IsEncrypted(t *testing.T) {
	p := new(PlainProtocol)
	assert.False(t, p.IsEncrypted())
}