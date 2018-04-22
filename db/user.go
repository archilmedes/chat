package db

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

func (u *User) AddFriend(displayName, macAddress, ipAddress, friendUsername string) bool {
	return addFriend(u.Username, displayName, macAddress, ipAddress, friendUsername)
}

func (u *User) GetFriends() []Friend {
	return getFriends(u.Username)
}

func (u *User) IsFriendsWith(displayName string) bool {
	return areFriends(u.Username, displayName)
}

func (u *User) GetFriendByDisplayName(friendDisplayName string) *Friend {
	return getFriendByDisplayName(u.Username, friendDisplayName)
}

func (u *User) GetFriendByUsernameAndMAC(friendUsername, friendMAC string) *Friend {
	return getFriendByUsernameAndMAC(u.Username, friendUsername, friendMAC)
}

func (u *User) UpdateMyIP() bool {
	return updateFriendIP(u.MAC, u.IP)
}

// Log in the user or return null.
func UserLogin(username string, password string) *User {
	return GetUser(username, password)
}

