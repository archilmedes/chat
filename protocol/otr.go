package protocol

import (
	"golang.org/x/crypto/otr"
	"crypto/rand"
	"log"
)

type OTRProtocol struct {
	Protocol
	converse *otr.Conversation
}

func NewOTRProtocol() *OTRProtocol {
	reader := rand.Reader
	privKey := new(otr.PrivateKey)
	privKey.Generate(reader)
	pubKey := new(otr.PublicKey)
	privKey.PublicKey = *pubKey
	converse := new(otr.Conversation)
	converse.PrivateKey = privKey
	return &OTRProtocol{converse: converse}
}

func (o OTRProtocol) Encrypt(in []byte) ([][]byte, error) {
	cipherText, err := o.converse.Send(in)
	if err != nil {
		log.Fatal(err)
	}
	return cipherText, nil
}

func (o OTRProtocol) Decrypt(in []byte) ([]byte, error) {
	out, encrypted, _, _, err := o.converse.Receive(in)
	if !encrypted {
		log.Fatalf("Message received was not encrypted")
	}
	if err != nil {
		log.Fatal(err)
	}
	return out, nil
}

func (o OTRProtocol) IsEncrypted() bool {
	return o.converse.IsEncrypted()
}

func (o OTRProtocol) EndSession() {
	o.converse.End()
}
