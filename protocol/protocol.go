// Contains protocol definition and plain protocol implementation
package protocol

type Protocol interface {
	Encrypt(in []byte) ([][]byte, error)
	Decrypt(cypher []byte) ([]byte, error)
	IsEncrypted() bool
	IsActive() bool
	NewSession() (string, error)
	EndSession()
	Serialize() []byte
}

// Type of protocol that just lets text pass through
type PlainProtocol struct {
	Protocol
}

// Encrypts the text by adding it into a 2D byte array
func (p PlainProtocol) Encrypt(in []byte) ([][]byte, error) {
	b := make([][]byte, 1)
	b[0] = in
	return b, nil
}

// Decrypts the message by just returning it
func (p PlainProtocol) Decrypt(cypher []byte) ([]byte, error) {
	return cypher, nil
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