package core

import (
	"log"
	"fmt"
	secure "chat/protocol"
)

type User struct {
	proto secure.Protocol // protocol used for encrypting messages
	login string // unique identifier for the user
}

func NewSecureUser(login string) *User {
	u := new(User)
	u.login = login
	u.proto = secure.NewOTRProtocol()
	return u
}

func (u *User) Persist() {
	protoBytes := u.proto.Serialize()
	login := u.login
	// TODO put method to save to DB, remove print
	fmt.Println("Saving " + string(protoBytes) + " login " + login)
}

func (u *User) Delete() {
	// TODO invoke delete by user login to DB
}

// Log in the user or return null. Return the user with the private key
func UserLogin(login string, password string) *User {
	// TODO Invoke DB
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