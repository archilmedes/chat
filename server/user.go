package server

import (
	"chat/db"
)

type User struct {
	IP, Username string
}

// Create a DatabaseUser to a core User object
func databaseUserToUser(dUser *db.DatabaseUser) (*User) {
	coreUser := new(User)
	coreUser.IP = (*dUser).IP
	coreUser.Username = (*dUser).Username
	return coreUser
}

// Persists the user to the database
func (u *User) Create(password string) (bool) {
	return db.AddUser(u.Username, password, u.IP)
}

// Delete the user from the database
func (u *User) Delete() (bool) {
	return db.DeleteUser(u.Username)
}

// Log in the user or return null.
func UserLogin(username string, password string) *User {
	dbUser := db.GetUser(username, password)
	return databaseUserToUser(dbUser)
}

func (u *User) SendMessage(toIpAddr string, text []byte) {

}

func (u *User) ReceiveMessage() {

}