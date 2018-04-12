// Contains base packages to implement the OTR protocol and encrypt messagse
package protocol

import (
	"golang.org/x/crypto/otr"
	"crypto/rand"
	"log"
)

type OTRProtocol struct {
	Protocol
	conv *otr.Conversation
	sess OTRSession
}

type OTRSession struct {
	SSID [8]byte
	fingerprint, privKey []byte
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
	return OTRProtocol{conv: conv}
}

func NewOTRProtocolFromKeys(privKeyBytes[]byte) OTRProtocol {
	privKey := new(otr.PrivateKey)
	privKey.Parse(privKeyBytes)
	conv := new(otr.Conversation)
	conv.PrivateKey = privKey
	conv.FragmentSize = fragmentSize
	return OTRProtocol{conv: conv}
}

// Encrypt the message
func (o OTRProtocol) Encrypt(in []byte) ([][]byte, error) {
	cipherText, err := o.conv.Send(in)
	if err != nil {
		log.Fatal(err)
	}
	return cipherText, nil
}

func (o OTRProtocol) Decrypt(in []byte) ([]byte, error) {
	out, encrypted, secChange, msgToPeer, err := o.conv.Receive(in)
	if err != nil {
		log.Fatal(err)
	}
	// Respond to handshake if handshake is established
	if len(msgToPeer) > 0 {
		log.Println("<OTR> Handshaking")
		for _, msg := range msgToPeer {
			// TODO server.sendMessage(destIp, msg), refactor to get a connection object here
			n := len(msg)
			if n < len(msg) {
				log.Panicln("<OTR> Handshake could not be established")
			}
			return msg, err
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
			o.conv.TheirPublicKey.Fingerprint(),
			o.conv.SSID,
		)
		// TODO Send OTRSession object to DB to save
		sess := new(OTRSession)
		sess.fingerprint = o.conv.TheirPublicKey.Fingerprint()
		sess.SSID = o.conv.SSID
		sess.privKey = o.conv.PrivateKey.Serialize(nil)
		o.sess = *sess
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
	return o.conv.IsEncrypted()
}

func (o OTRProtocol) EndSession() {
	o.conv.End()
}

func (o OTRProtocol) serialize() []byte {
	return o.conv.PrivateKey.Serialize(nil)
}