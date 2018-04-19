package protocol

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestPlainProtocol_IsActive(t *testing.T) {
	p := new(PlainProtocol)
	assert.True(t, p.IsActive())
}

func TestPlainProtocol_ToType(t *testing.T) {
	p := new(PlainProtocol)
	assert.Equal(t, PlainType, p.ToType())
}

func TestPlainProtocol_NewSession(t *testing.T) {
	p := new(PlainProtocol)
	firstMessage, err := p.NewSession()
	assert.Nil(t, err)
	assert.Equal(t, "", firstMessage)
}

func TestPlainProtocol_Serialize(t *testing.T) {
	exp := new(PlainProtocol)
	actual := CreateProtocolFromType(exp.ToType())
	actual.InitFromBytes(exp.Serialize())
	assert.Equal(t, *exp, actual)
}

func TestPlainProtocol_Serialize_NewSession(t *testing.T) {
	exp := new(PlainProtocol)
	exp.NewSession()
	actual := CreateProtocolFromType(exp.ToType())
	actual.InitFromBytes(exp.Serialize())
	assert.Equal(t, *exp, actual)
	assert.Equal(t, exp.SessionID, actual.GetSessionID())
}

func TestCreateProtocolFromType_plain(t *testing.T) {
	p := new(PlainProtocol)
	assert.Equal(t, PlainProtocol{}, CreateProtocolFromType(p.ToType()))
}
