// Package protocol contains basic protocol functionality to encrypt, decrypt messages with
// plaintext and OTR protocols
package protocol

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
)

// A generic Protocol interface to handle common protocol methods
type Protocol interface {
	InitFromBytes([]byte) error
	Encrypt(in []byte) ([][]byte, error)
	Decrypt(cypher []byte) ([][]byte, error)
	IsEncrypted() bool
	IsActive() bool
	NewSession() (string, error)
	GetSessionID() uint64
	EndSession()
	Serialize() []byte
	ToType() string
}

// Type of protocol that just lets text pass through without applying any encryption
type PlainProtocol struct {
	Protocol
	SessionID uint64
}

const (
	PlainType = "plain"
	OTRType   = "otr"
)

// Wrap the a byte array into an array of messages
func wrapMessage(in []byte) [][]byte {
	b := make([][]byte, 1)
	b[0] = in
	return b
}

func init() {
	gob.Register(&PlainProtocol{})
	gob.Register(&OTRProtocol{})
}

// Given the protocol type, reconstruct the subclass
// Should be used in conjunction with Protocol#InitFromBytes to re-create a protocol
func CreateProtocolFromType(protoType string) Protocol {
	if protoType == PlainType {
		return new(PlainProtocol)
	} else if protoType == OTRType {
		return NewOTRProtocol()
	} else {
		panic(errors.New(fmt.Sprintf("CreateProtocolFromType: %s", protoType)))
	}
}

func (p *PlainProtocol) InitFromBytes(dec []byte) error {
	decBuf := bytes.NewBuffer(dec)
	return gob.NewDecoder(decBuf).Decode(p)
}

// Encrypts the text by adding it into a 2D byte array
func (p *PlainProtocol) Encrypt(in []byte) ([][]byte, error) {
	return wrapMessage(in), nil
}

// Decrypts the message by just returning it
func (p *PlainProtocol) Decrypt(dec []byte) ([][]byte, error) {
	return wrapMessage(dec), nil
}

// Always returns false as a plain protocol is never encrypted
func (p *PlainProtocol) IsEncrypted() bool {
	return false
}

// Always returns true as a plain protocol is always active
func (p *PlainProtocol) IsActive() bool {
	return true
}

// Start a new plain protocol session
func (p *PlainProtocol) NewSession() (string, error) {
	var n uint64
	err := binary.Read(rand.Reader, binary.LittleEndian, &n)
	p.SessionID = n
	return "", err
}

func (p *PlainProtocol) GetSessionID() uint64 {
	return p.SessionID
}

// Ends a plain protocol session
func (p *PlainProtocol) EndSession() {
	p.SessionID = 0
}

// Serialize the protocol to save in a database
func (p *PlainProtocol) Serialize() []byte {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(p)
	return b.Bytes()
}

// Converts plain protocol to type
func (p *PlainProtocol) ToType() string {
	return PlainType
}
