package protocol

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/otr"
	"testing"
)

func TestOTRProtocol_EndSession(t *testing.T) {
	o := NewOTRProtocol()
	o.EndSession()
	assert.False(t, o.Conv.IsEncrypted())
}

func TestOTRProtocol_Handshake(t *testing.T) {
	alice, bob := NewOTRProtocol(), NewOTRProtocol()
	createOTRLink(t, alice, bob)
}

// Inspired by official OTR tests in Golang here: https://github.com/keybase/go-crypto
func createOTRLink(t *testing.T, alice *OTRProtocol, bob *OTRProtocol) {
	var aMsg, bMsg [][]byte
	aMsg = append(aMsg, []byte(otr.QueryMessage))
	// Simulate a handshake by just sending messages between two users
	for len(aMsg[0]) > 0 || len(bMsg[0]) > 0 {
		var err error
		bMsg = [][]byte{}
		for _, msg := range aMsg {
			bMsg, err = bob.Decrypt(msg)
			assert.Error(t, OTRHandshakeStep{}, err)
		}
		aMsg = [][]byte{}
		for _, msg := range bMsg {
			aMsg, err = alice.Decrypt(msg)
			assert.Error(t, OTRHandshakeStep{}, err)
		}
	}
	assert.True(t, alice.IsEncrypted())
	assert.True(t, bob.IsEncrypted())
	assert.True(t, alice.IsActive())
	assert.True(t, bob.IsActive())
	assert.Equal(t, alice.GetSessionID(), bob.GetSessionID())
}

func TestOTRProtocol_SendAndReceiveMessages(t *testing.T) {
	andrew, sameet := NewOTRProtocol(), NewOTRProtocol()
	createOTRLink(t, andrew, sameet)

	var testMessage = "Want to play fortnite?"
	cyp, err := andrew.Encrypt([]byte(testMessage))
	assert.Nil(t, err)

	for _, msg := range cyp {
		out, err := sameet.Decrypt(msg)
		assert.Nil(t, err)
		assert.Equal(t, testMessage, string(out[0]), "Strings should be equivalent")
	}
}

func TestOTRProtocol_NewSession(t *testing.T) {
	proto := new(OTRProtocol)
	firstMessage, err := proto.NewSession()
	assert.Nil(t, err)
	assert.Equal(t, otr.QueryMessage, firstMessage)
}

func TestOTRProtocol_Serialize_BeforeHandshake(t *testing.T) {
	exp := NewOTRProtocol()

	ott := CreateProtocolFromType(exp.ToType())
	ott.InitFromBytes(exp.Conv.PrivateKey.Serialize(nil))
	assert.Equal(t, *exp, *ott.(*OTRProtocol))
}

func TestOTRProtocol_Serialize_AfterHandshake(t *testing.T) {
	alice, bob := NewOTRProtocol(), NewOTRProtocol()
	createOTRLink(t, alice, bob)

	ott := CreateProtocolFromType(alice.ToType())
	ott.InitFromBytes(alice.Conv.PrivateKey.Serialize(nil))
	assert.Equal(t, *alice.Conv.PrivateKey, *ott.(*OTRProtocol).Conv.PrivateKey)
}

func TestCreateProtocolFromType_otr(t *testing.T) {
	p := new(OTRProtocol)
	proto := CreateProtocolFromType(p.ToType())
	assert.NotNil(t, proto)
	_, ok := proto.(*OTRProtocol)
	assert.True(t, ok)
}
