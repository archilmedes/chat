// Contains protocol definition and plain protocol implementation
package protocol

import (
	"errors"
	"fmt"
)

type Protocol interface {
	Encrypt(in []byte) ([][]byte, error)
	Decrypt(cypher []byte) ([][]byte, error)
	IsEncrypted() bool
	IsActive() bool
	NewSession() (string, error)
	EndSession()
	Serialize() []byte
	ToType() string
}

const (
	PlainType = "plain"
	OTRType = "otr"
)

// Type of protocol that just lets text pass through
type PlainProtocol struct {
	Protocol
}

func wrapMessage(in[] byte) ([][]byte) {
	b := make([][]byte, 1)
	b[0] = in
	return b
}

func CreateProtocolFromType(protoType string) Protocol {
	if protoType == PlainType {
		return PlainProtocol{}
	} else if protoType == OTRType {
		return NewOTRProtocol()
	} else {
		panic(errors.New(fmt.Sprintf("CreateProtocolFromType: %s", protoType)))
	}
}

// Encrypts the text by adding it into a 2D byte array
func (p PlainProtocol) Encrypt(in []byte) ([][]byte, error) {
	return wrapMessage(in), nil
}
// Decrypts the message by just returning it
func (p PlainProtocol) Decrypt(dec []byte) ([][]byte, error) {
	return wrapMessage(dec), nil
}

// Always returns false as a plain protocol is never encrypted
func (p PlainProtocol) IsEncrypted() bool {
	return false
}

// Always returns true as a plain protocol is always active
func (p PlainProtocol) IsActive() bool {
	return true
}

// Start a new plain protocol session
func (p PlainProtocol) NewSession() (string, error) {
	return "", nil
}

// Ends a plain protocol session
func (p PlainProtocol) EndSession() {
	// no-op
}

// Serialize the protocol to save in a database
func (p PlainProtocol) Serialize() []byte {
	return []byte(nil)
}

// Converts plain protocol to type
func (p PlainProtocol) ToType() string {
	return PlainType
}