package db

import "github.com/wavyllama/chat/config"

// Stores a user's information
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

func (u *User) GetFriend(friendUsername, friendMAC string) *Friend {
	return db.GetFriend(u.Username, friendUsername, friendMAC)
}

// Log in the user or return null.
func UserLogin(username string, password string) *User {
	return GetUser(username, password)
}

