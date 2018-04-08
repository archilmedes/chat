package protocol

import "log"

type User struct {
	proto Protocol // protocol used for encrypting messages
	login string // unique identifier for the user
}

func NewSecureUser() *User {
	// TODO Week 2/3: create a unique login string of length n
	// that users can use to login and restore conversation
	login := "test"
	u := &User{login: login}
	u.proto = NewOTRProtocol()
	return u
}

// Log in the user or return null. Return the user with the private key
func UserLogin(login string, password string) *User {
	// TODO Implement in Week 3
	return nil
}

func (u *User) NewSession(destIp string) (bool, error) {
	// TODO send otr.QueryMessage to initiate OTR handshake
	return true, nil
}

func (u *User) ReceiveMessage(enc []byte) ([]byte, error) {
	msg, err := u.proto.Decrypt(enc)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return msg, nil
}

func (u *User) SendMessage(msg[] byte) (bool, error) {
	_, err := u.proto.Encrypt(msg)
	if err != nil {
		return false, err
	}
	return false, nil
}

func (u *User) EndSession() {
	u.proto.EndSession()
}

func (u *User) IsSecure() bool {
	return u.proto.IsEncrypted()
}