// Contains protocol definition and plain protocol implementation
package protocol

type Protocol interface {
	Encrypt(in []byte) ([][]byte, error)
	Decrypt(cypher []byte) ([]byte, error)
	IsEncrypted() bool
	EndSession()
	Serialize() []byte
}

type PlainProtocol struct {
	Protocol
}

func (p PlainProtocol) Encrypt(in []byte) ([][]byte, error) {
	b := make([][]byte, 1)
	b[0] = in
	return b, nil
}

func (p PlainProtocol) Decrypt(cypher []byte) ([]byte, error) {
	return cypher, nil
}

func (p PlainProtocol) IsEncrypted() bool {
	return false
}

func (p PlainProtocol) EndSession() {
	// no-op
}

func (p PlainProtocol) Serialize() []byte {
	return []byte(nil)
}