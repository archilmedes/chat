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

func (o OTRProtocol) Decrypt(in []byte) ([]byte, error) {
	out, encrypted, secChange, msgToPeer, err := o.Conv.Receive(in)
	if err != nil {
		log.Fatal(err)
	}
	// Respond to handshake if handshake is established
	if len(msgToPeer) > 0 {
		log.Println("<OTR> Handshaking")
		for _, msg := range msgToPeer {
			return msg, OTRHandshakeStep{}
		}
	}
	switch secChange {
	case otr.NoChange:
		// If it's encrypted, just return the decrypted message out
		if encrypted && o.IsEncrypted() {
			return out, nil
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
		return out, nil
	case otr.ConversationEnded:
		log.Println("<OTR> Conversation ended")
		o.EndSession()
		return out, nil
	default:
		log.Fatal("<OTR> SMP not implemented")
	}

	return out, nil
}

func (o OTRProtocol) IsEncrypted() bool {
	return o.Conv.IsEncrypted()
}

func (o OTRProtocol) EndSession() {
	o.Conv.End()
}

func (o OTRProtocol) serialize() []byte {
	return o.Conv.PrivateKey.Serialize(nil)
}