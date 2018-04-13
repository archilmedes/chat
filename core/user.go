package core

import (
	"fmt"
	secure "chat/protocol"
)

type User struct {
	Proto secure.Protocol // protocol used for encrypting messages
	IP, Login string
}

func NewSecureUser(login string) *User {
	u := new(User)
	u.Login = login
	u.Proto = secure.NewOTRProtocol()
	return u
}

func (u *User) Persist() {
	protoBytes := u.Proto.Serialize()
	login := u.Login
	// TODO put method to save to DB, remove print
	fmt.Println("Saving " + string(protoBytes) + " Login " + login)
}

func (u *User) Delete() {
	// TODO invoke delete by user Login to DB
}

// Log in the user or return null. Return the user with the private key
func UserLogin(login string, password string) *User {
	// TODO Invoke DB
	return nil
}

// Receives the message from a user
func (u *User) ReceiveMessage(from *User, enc []byte) ([]byte, error) {
	return u.Proto.Decrypt(enc)
}

func (u *User) EncryptMessage(msg string) ([][]byte, error) {
	cypher, err := u.Proto.Encrypt([]byte(msg))
	if err != nil {
		return nil, err
	}
	return cypher, nil
}

func (u *User) EndSession() {
	u.Proto.EndSession()
}

func (u *User) IsSecure() bool {
	return u.Proto.IsEncrypted()
}