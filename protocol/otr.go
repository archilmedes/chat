package protocol

import (
	"crypto/rand"
	"golang.org/x/crypto/otr"
	"log"
	"fmt"
	"strconv"
	"encoding/gob"
	"bytes"
	"math/big"
	"io"
	"crypto/sha256"
)

// An OTR protocol that contains the crypto/otr struct
type OTRProtocol struct {
	Protocol
	Conv    *otr.Conversation
}

type OTRProtocolSerialization struct {
	// PrivateKey contains the private key to use to sign key exchanges.
	PrivateKey *otr.PrivateKey

	// Rand can be set to override the entropy source. Otherwise,
	// crypto/rand will be used.
	Rand io.Reader
	// If FragmentSize is set, all messages produced by Receive and Send
	// will be fragmented into messages of, at most, this number of bytes.
	FragmentSize int

	// Once Receive has returned NewKeys once, the following fields are
	// valid.
	SSID           [8]byte
	TheirPublicKey otr.PublicKey

	state, authState int

	r       [16]byte
	x, y    *big.Int
	gx, gy  *big.Int
	gxBytes []byte
	digest  [sha256.Size]byte

	//revealKeys, sigKeys akeKeys
	//c      [16]byte
	//m1, m2 [32]byte

	myKeyId         uint32
	myCurrentDHPub  *big.Int
	myCurrentDHPriv *big.Int
	myLastDHPub     *big.Int
	myLastDHPriv    *big.Int

	theirKeyId        uint32
	theirCurrentDHPub *big.Int
	theirLastDHPub    *big.Int

	//keySlots [4]keySlot

	myCounter    [8]byte
	theirLastCtr [8]byte
	oldMACs      []byte

	k, n int // fragment state
	frag []byte

	//smp smpState
}

// A type of error that indicates that the protocol is undergoing the handshake
type OTRHandshakeStep struct {
	error
}

const (
	fragmentSize = 1000
)

// Create a new OTR session, with new keys and a new Conversation
func NewOTRProtocol() *OTRProtocol {
	privKey := new(otr.PrivateKey)
	privKey.Generate(rand.Reader)
	conv := new(otr.Conversation)
	conv.PrivateKey = privKey
	conv.FragmentSize = fragmentSize
	return &OTRProtocol{Conv: conv}
}

func (o *OTRProtocol) InitFromBytes(dec []byte) error {
	decBuf := bytes.NewBuffer(dec)
	return gob.NewDecoder(decBuf).Decode(o)
}

// Encrypt the message
func (o *OTRProtocol) Encrypt(in []byte) ([][]byte, error) {
	cipherText, err := o.Conv.Send(in)
	if err != nil {
		log.Fatal(err)
	}
	return cipherText, nil
}

// Decrypt the message and handle OTR protocol
func (o *OTRProtocol) Decrypt(in []byte) ([][]byte, error) {
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
func (o *OTRProtocol) NewSession() (string, error) {
	return otr.QueryMessage, nil
}

// Get the SessionID for an OTR session if it exists, or 0 if no session id exists
func (o *OTRProtocol) GetSessionID() uint64 {
	if !o.IsActive() {
		return 0
	}
	// Convert a [8]byte -> string -> uint64
	SSID := fmt.Sprintf("%x", o.Conv.SSID)
	sessionId, err := strconv.ParseUint(SSID, 16, 64)
	if err != nil {
		fmt.Errorf("Error getting session id: %s\n", err.Error())
	}
	return sessionId
}

// Returns true if an OTR conversation is now encrypted
func (o *OTRProtocol) IsEncrypted() bool {
	return o.Conv.IsEncrypted()
}

// Returns true if an active OTR session has been created and in use
func (o *OTRProtocol) IsActive() bool {
	return o.IsEncrypted()
}

// Ends the OTR conversation
func (o *OTRProtocol) EndSession() {
	o.Conv.End()
}

// Serialize the private key of an OTR protocol, which is all that is needed to recreate
func (o *OTRProtocol) Serialize() []byte {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(o)
	return b.Bytes()
}

// Return type of protocol
func (o *OTRProtocol) ToType() string {
	return OTRType
}
