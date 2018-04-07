package protocol

type User struct {
	proto Protocol // protocol used for encrypting messages
	login string // unique identifier for the user
}

func NewSecureUser() *User {
	login := "test" // TODO create a unique login string of length n
	u := &User{login: login}
	u.proto = NewOTRProtocol()
	return u
}

// Log in the user or return null. Return the user with the private key
func (u *User) Login(login string, password string) *User {
	// TODO Implement in Week 3
	return nil
}

func (u *User) NewSession() {
	u.proto = NewOTRProtocol()
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