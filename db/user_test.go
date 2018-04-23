package db

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const(
	mac = "AA:BB:CC:DD:EE"
	ip = "1.2.3.4"
	username = "sameet"
	friendDisplayName = "Bob"
)

func getFakeUser() *User {
	SetupEmptyTestDatabase()
	user := NewUser(mac, ip, username)
	// Added this line to simulate how every user is friends with themselves
	user.AddFriend(Self, user.MAC, user.IP, user.Username)
	return user
}

func TestUser_AddFriend(t *testing.T) {
	user := getFakeUser()

	assert.True(t, user.AddFriend(friendDisplayName, "doesntmatter", "5.6.7.8", "bobbyB"))
	friend := user.GetFriendByDisplayName(friendDisplayName)
	assert.NotNil(t, friend)

	assert.Equal(t, friendDisplayName, friend.DisplayName)

	assert.True(t, user.IsFriendsWith(friendDisplayName))
	assert.False(t, user.IsFriendsWith("Charlie"))
}

func TestUser_GetConversationHistory(t *testing.T) {
	user := getFakeUser()

	assert.True(t, user.AddFriend(friendDisplayName, "doesntmatter", "5.6.7.8", "bobbyB"))
}

func TestUser_IsFriendOnline(t *testing.T) {
	user := getFakeUser()

	online, _ := user.IsFriendOnline(friendDisplayName)
	assert.False(t, online)

	online, _ = user.IsFriendOnline(Self)
	assert.True(t, online)
}

func TestUser_UpdateMyIP(t *testing.T) {
	user := getFakeUser()
	user.IP = "newip"
	user.UpdateMyIP()

	myself := user.GetFriends()[0]
	assert.Equal(t, user.IP, myself.IP)
}

func TestUser_GetSessions(t *testing.T) {
	user := getFakeUser()
	assert.Equal(t, 0, len(user.GetSessions(Self)))
}