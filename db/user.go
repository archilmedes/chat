package db

type User struct {
	Username, MAC, IP string
}

// Persists the user to the database
func (u *User) Create(password string) bool {
	return AddUser(u.Username, password, u.IP)
}

// Delete the user from the database
func (u *User) Delete() bool {
	return DeleteUser(u.Username)
}

// Log in the user or return null.
func UserLogin(username string, password string) *User {
	return GetUser(username, password)
}

func (u *User) SendMessage(toIpAddr string, text []byte) {

}

func (u *User) ReceiveMessage() {

}
