package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFriends(t *testing.T) {
	SetupDatabaseForTests(t)
	FriendSetup(t)
	AreFriendsTest(t)
	AddUpdateDeleteFriendTest(t)
}

func FriendSetup(t *testing.T) {
	friends := getFriends("karateAMD")
	assert.Equal(t, 3, len(friends))
	assert.Equal(t, "archilmedes", friends[0].Username)
	assert.Equal(t, "10.192.345.987", friends[1].IP)
	assert.Equal(t, "11:11:11:11", friends[2].MAC)
	assert.Equal(t, "andrew", friends[2].DisplayName)
}

func AreFriendsTest(t *testing.T) {
	assert.True(t, areFriends("karateAMD", "sameet"))
	assert.False(t, areFriends("andrew", "archilmedes"))
}

func AddUpdateDeleteFriendTest(t *testing.T) {
	assert.True(t, addFriend("alice123", "andrew", "11:11:11:11", "987.654.321", "andrew"))
	friends := getFriends("alice123")
	assert.Equal(t, 2, len(friends))
	assert.True(t, updateFriendIP("alice123", "11:11:11:11", "444.444.444.444"))
	friends = getFriends("alice123")
	assert.Equal(t, "444.444.444.444", friends[1].IP)
	assert.True(t, deleteFriend("alice123", "andrew"))
	friends = getFriends("alice123")
	assert.Equal(t, 1, len(friends))
}
