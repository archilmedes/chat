package db

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func FriendsTest(t *testing.T) {
	FriendSetup(t)
	AreFriendsTest(t)
	AddUpdateDeleteFriendTest(t)
}

func FriendSetup(t *testing.T) {
	friends := GetFriends("karateAMD")
	assert.Equal(t, 3, len(friends))
	assert.Equal(t, "archilmedes", friends[0].Username)
	assert.Equal(t, "10.192.345.987", friends[1].IP)
	assert.Equal(t,"11:11:11:11", friends[2].MAC)
	assert.Equal(t, "andrew", friends[2].DisplayName)
}

func AreFriendsTest(t *testing.T) {
	assert.True(t, AreFriendsName("karateAMD", "sameet"))
	assert.False(t, AreFriendsName("andrew", "archilmedes"))
	assert.True(t, AreFriendsMac("archilmedes", "01:23:45:67"))
	assert.False(t, AreFriendsMac("karateAMD", "ff:ff:ff:ff"))
}

func AddUpdateDeleteFriendTest(t *testing.T) {
	assert.True(t, AddFriend("alice123", "andrew", "11:11:11:11", "987.654.321", "andrew"))
	friends := GetFriends("alice123")
	assert.Equal(t, 2, len(friends))
	assert.True(t, UpdateFriendIp("11:11:11:11", "444.444.444.444"))
	friends = GetFriends("alice123")
	assert.Equal(t, "444.444.444.444", friends[1].IP)
	assert.True(t, DeleteFriend("alice123", "11:11:11:11"))
	friends = GetFriends("alice123")
	assert.Equal(t, 1, len(friends))
}