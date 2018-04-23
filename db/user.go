package db

import (
	"log"
	"github.com/wavyllama/chat/protocol"
	"os/exec"
	"strings"
	"time"
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

func (u *User) GetSessions(friendDisplayName string) []Session {
	return getUserSessions(u.Username)
}

func (u *User) UpdateMyIP() bool {
	return updateFriendIP(u.MAC, u.IP)
}

// Checks if a friend is online, and return a timestamp of when they were last online
func (u *User) IsFriendOnline(friendDisplayName string) (bool, time.Time) {
	var lastSeenTime time.Time
	friend := u.GetFriendByDisplayName(friendDisplayName)
	sessions := u.GetSessions(friendDisplayName)
	if friend == nil || len(sessions) == 0 {
		return false, lastSeenTime
	}
	out, _ := exec.Command("ping", friend.IP, "-c 5", "-i 3", "-w 10").Output()
	friendOnline := !strings.Contains(string(out), "Destination Host Unreachable")

	// If the friend is online now, then they are available now
	if friendOnline {
		lastSeenTime = time.Now()
	} else {
		// Otherwise their last message in the last session is when they were last online
		//lastSeenTime, _ = time.Parse(time.RFC3339, sessions[len(sessions) - 1].timestamp)
		messages := getSessionMessages(sessions[len(sessions) - 1].SSID)
		for i := len(messages) - 1; i >= 0; i-- {
			if messages[i].SentOrReceived == 1 {
				lastSeenTime, _ = time.Parse(time.RFC3339, messages[i].Timestamp)
			}
		}
	}
	return friendOnline, lastSeenTime
}

// Fetch conversations between another friend and decrypts the contents of the messages everything
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