package db

import (
	"os/exec"
	"strings"
	"time"
)

// Stores a user's information
type User struct {
	MAC, IP, Username string
}

func NewUser(mac, ip, username string) *User {
	return &User{MAC: mac, IP: ip, Username: username}
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

func (u *User) DeleteFriend(displayName string) bool {
	return deleteFriend(u.Username, displayName)
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

// Return a user's sessions in descending order timestamp (most recent session first)
func (u *User) GetSessions(friendDisplayName string) []Session {
	return getUserSessions(u.Username)
}

func (u *User) UpdateMyIP() bool {
	return updateFriendIP(u.Username, u.MAC, u.IP)
}

// Checks if a friend is online, and return a timestamp of when they were last online
func (u *User) IsFriendOnline(friendDisplayName string) (bool, time.Time) {
	friend := u.GetFriendByDisplayName(friendDisplayName)
	// If they aren't a friend or you've never communicated with him/her
	if friend == nil {
		return false, time.Time{}
	}
	out, _ := exec.Command("ping", friend.IP, "-c 5", "-i 3", "-w 10").Output()
	friendOnline := !strings.Contains(string(out), "Destination Host Unreachable")
	// If the friend is online now, then they are available now
	if friendOnline || friend.DisplayName == Self {
		return true, time.Now()
	}

	sessions := u.GetSessions(friendDisplayName)
	var lastSeenTime time.Time
	// Otherwise their last message in the last session is when they were last online
	messages := getSessionMessages(sessions[len(sessions)-1].SSID)
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].SentOrReceived == Received {
			lastSeenTime, _ = time.Parse(time.RFC3339, messages[i].Timestamp)
		}
	}
	return friendOnline, lastSeenTime
}

// Fetch conversations between another friend and decrypts the contents of the messages everything
func (u *User) GetConversationHistory(friendDisplayName string) []Conversation {
	return getConversationsWithFriend(u.Username, friendDisplayName)
}
