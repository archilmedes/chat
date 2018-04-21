package db

import (
	"log"
	"github.com/wavyllama/chat/protocol"
)

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

// Fetch conversations between another friend
func (u *User) GetConversationHistory(friendDisplayName string) [][]byte {
	converse := GetConversationUsers(u.Username, friendDisplayName)
	var messages [][]byte
	for _, c := range converse {
		proto := protocol.CreateProtocolFromType(c.Session.ProtocolType)
		proto.InitFromBytes([]byte(c.Session.ProtocolValue))
		dec, err := proto.Decrypt([]byte(c.Message.Text))
		if err != nil {
			log.Printf("GetMessages: %s", err.Error())
		}
		messages = append(messages, dec[0])
	}
	return messages
}