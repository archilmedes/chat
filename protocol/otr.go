// Package protocol contains basic protocol functionality to encrypt, decrypt messages with
// plaintext and OTR protocols
package protocol

import (
	"crypto/rand"
	"golang.org/x/crypto/otr"
	"log"
	"fmt"
	"strconv"
)

type OTRProtocol struct {
	Protocol
	Conv    *otr.Conversation
}

type OTRHandshakeStep struct {
	error
}

const (
	fragmentSize = 1000
)

// Create a new OTR session, with new keys and a new Conversation
func NewOTRProtocol() OTRProtocol {
	privKey := new(otr.PrivateKey)
	privKey.Generate(rand.Reader)
	conv := new(otr.Conversation)
	conv.PrivateKey = privKey
	conv.FragmentSize = fragmentSize
	return OTRProtocol{Conv: conv}
}

func (o OTRProtocol) InitFromBytes(privKeyBytes []byte) error {
	privKey := new(otr.PrivateKey)
	privKey.Parse(privKeyBytes)
	o.Conv = new(otr.Conversation)
	o.Conv.PrivateKey = privKey
	o.Conv.FragmentSize = fragmentSize
	return nil
}

// Encrypt the message
func (o OTRProtocol) Encrypt(in []byte) ([][]byte, error) {
	cipherText, err := o.Conv.Send(in)
	if err != nil {
		log.Fatal(err)
	}
	return cipherText, nil
}

// Decrypt the message and handle OTR protocol
func (o OTRProtocol) Decrypt(in []byte) ([][]byte, error) {
	out, encrypted, secChange, msgToPeer, err := o.Conv.Receive(in)
	if err != nil {
		log.Fatal(err)
	}
	// Respond to handshake if handshake is established with the message to send back
	if len(msgToPeer) > 0 {
		return msgToPeer, OTRHandshakeStep{}
	}
	switch secChange {
	case otr.NoChange:
		// If it's encrypted, just return the decrypted message out
		if encrypted && o.IsEncrypted() {
			return wrapMessage(out), nil
		}
	case otr.NewKeys:
		log.Println("<OTR> Key exchange completed. You are now in a secure session.")
		return wrapMessage(out), nil
	case otr.ConversationEnded:
		o.EndSession()
		return wrapMessage(out), nil
	default:
		log.Fatal("<OTR> SMP not implemented")
	}

	return wrapMessage(out), nil
}

// Create a new session
func (o OTRProtocol) NewSession() (string, error) {
	return otr.QueryMessage, nil
}

func (o OTRProtocol) GetSessionID() uint64 {
	SSID := fmt.Sprintf("%x", o.Conv.SSID)
	sessionId, err :=strconv.ParseUint(SSID, 16, 64)
	if err != nil {
		fmt.Errorf("Error getting session id: %s\n", err.Error())
	}
	return sessionId
}

// Returns true if an OTR conversation is now encrypted
func (o OTRProtocol) IsEncrypted() bool {
	return o.Conv.IsEncrypted()
}

// Returns true if an OTR session has been created
func (o OTRProtocol) IsActive() bool {
	return o.IsEncrypted()
}

// Ends the OTR conversation
func (o OTRProtocol) EndSession() {
	o.Conv.End()
}

// Serialize the entire OTRProtocol object
func (o OTRProtocol) Serialize() []byte {
	//var b bytes.Buffer
	//gob.NewEncoder(&b).Encode(&o)
	//return b.Bytes()
	return o.Conv.PrivateKey.Serialize(nil)
}

// Return type of protocol
func (o OTRProtocol) ToType() string {
	return OTRType
}
