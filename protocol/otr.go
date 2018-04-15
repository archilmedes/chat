// Contains base packages to implement the OTR protocol and encrypt messagse
package protocol

import (
	"golang.org/x/crypto/otr"
	"crypto/rand"
	"log"
)

type OTRProtocol struct {
	Protocol
	Conv    *otr.Conversation
	Session OTRSession
}

type OTRSession struct {
	SSID                    [8]byte
	Fingerprint, PrivateKey []byte
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

// Recreates an OTR protocol given its private key
func NewOTRProtocolFromKeys(privKeyBytes[]byte) OTRProtocol {
	privKey := new(otr.PrivateKey)
	privKey.Parse(privKeyBytes)
	conv := new(otr.Conversation)
	conv.PrivateKey = privKey
	conv.FragmentSize = fragmentSize
	return OTRProtocol{Conv: conv}
}

// Encrypt the message
func (o OTRProtocol) Encrypt(in []byte) ([][]byte, error) {
	cipherText, err := o.Conv.Send(in)
	if err != nil {
		log.Fatal(err)
	}
	return cipherText, nil
}

func (o OTRProtocol) Decrypt(in []byte) ([][]byte, error) {
	out, encrypted, secChange, msgToPeer, err := o.Conv.Receive(in)
	if err != nil {
		log.Fatal(err)
	}
	// Respond to handshake if handshake is established
	if len(msgToPeer) > 0 {
		log.Println("<OTR> Handshaking")
		return msgToPeer, OTRHandshakeStep{}
	}
	switch secChange {
	case otr.NoChange:
		// If it's encrypted, just return the decrypted message out
		if encrypted && o.IsEncrypted() {
			return wrapMessage(out), nil
		}
	case otr.NewKeys:
		log.Printf("<OTR> Key exchange completed.\nFingerprint:%x\nSSID:%x\n",
			o.Conv.TheirPublicKey.Fingerprint(),
			o.Conv.SSID,
		)
		// TODO Send OTRSession object to DB to save
		sess := new(OTRSession)
		sess.Fingerprint = o.Conv.TheirPublicKey.Fingerprint()
		sess.SSID = o.Conv.SSID
		sess.PrivateKey = o.Conv.PrivateKey.Serialize(nil)
		o.Session = *sess
		return wrapMessage(out), nil
	case otr.ConversationEnded:
		log.Println("<OTR> Conversation ended")
		o.EndSession()
		return wrapMessage(out), nil
	default:
		log.Fatal("<OTR> SMP not implemented")
	}

	return wrapMessage(out), nil
}

func (o OTRProtocol) NewSession() (string, error) {
	return otr.QueryMessage, nil
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

func (o OTRProtocol) Serialize() []byte {
	return o.Conv.PrivateKey.Serialize(nil)
}

func (o OTRProtocol) ToType() string {
	return OTRType
}