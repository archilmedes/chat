package protocol

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/otr"
)

func TestOTRProtocol_Encrypt(t *testing.T) {
	proto := NewOTRProtocol()
	// TODO mock or something
	cipherText, _ := proto.Encrypt([]byte("message"))
	fmt.Println(cipherText)
}

func TestOTRProtocol_EndSession(t *testing.T) {
	o := NewOTRProtocol()
	o.EndSession()
	assert.False(t, o.conv.IsEncrypted())
}

// Inspired by official OTR tests in Golang here: https://github.com/keybase/go-crypto
func TestOTRProtocol_Handshake(t *testing.T) {
	alice, bob := NewOTRProtocol(), NewOTRProtocol()

	var aMsg, bMsg [][]byte
	var out []byte
	var err error
	aMsg = append(aMsg, []byte(otr.QueryMessage))
	var aSecChange, bSecChange otr.SecurityChange

	// Simulate a handshake by just sending messages between two users
	for ; len(aMsg) > 0 || len(bMsg) > 0; {
		bMsg = nil
		for _, msg := range aMsg {
			out, _, bSecChange, bMsg, err = bob.conv.Receive(msg)
			assert.Len(t, out, 0, "Should not generate output during key exchange")
			assert.Nil(t, err, "Error message not nil: %s", err)
		}

		aMsg = nil
		for _, msg := range bMsg {
			out, _, aSecChange, aMsg, err = alice.conv.Receive(msg)
			assert.Len(t, out, 0, "Should not generate output during key exchange")
			assert.Nil(t, err, "Error message not nil: %s", err)
		}
	}

	assert.Equal(t, otr.NewKeys, aSecChange, "Alice should have signaled NewKeys")
	assert.Equal(t, otr.NewKeys, bSecChange, "Bob should have signaled NewKeys")

	assert.True(t, alice.IsEncrypted(), "Alice should be encrypted")
	assert.True(t, bob.IsEncrypted(), "Alice should be encrypted")
}