package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsers(t *testing.T) {
	SetupTestDatabase()
	users := QueryUsers()
	UserSetupTest(t, users)
	UserExistsTest(t)
	InsertUserTest(t)
	UpdateUserTest(t)
	ValidateCredentialsTest(t)
}

func UpdateUserTest(t *testing.T) {
	assert.True(t, UpdateUserIP("sameetandpotatoes", "333.333.333.333"))
	assert.True(t, UpdateUserPassword("sameetandpotatoes", "p0t8t035AreCool"))
	users := QueryUsers()
	assert.Equal(t, 7, len(users))
	assert.Equal(t, "andrew", users[5].Username)
	assert.Equal(t, "888.888.888", users[5].IP)
}

func InsertUserTest(t *testing.T) {
	assert.True(t, AddUser("tempUser", "tempPass", "666.666.666.666"))
	users := QueryUsers()
	assert.Equal(t, 7, len(users))
}

func ValidateCredentialsTest(t *testing.T) {
	assert.NotNil(t, GetUser("karateAMD", "pwd123"))
	assert.Nil(t, GetUser("Sameet", "iLuvMacs"))
	assert.Nil(t, GetUser("sameetandpotatoes", "linuxFTW"))
}

func UserExistsTest(t *testing.T) {
	assert.True(t, UserExists("sameetandpotatoes"))
	assert.False(t, UserExists("fakeUser"))
}

func UserSetupTest(t *testing.T, users []User) {
	assert.Equal(t, 6, len(users))
	assert.Equal(t, "bob", users[1].Username)
	assert.Equal(t, "192.168.10.123", users[2].IP)
	assert.Equal(t, "10.192.345.987", users[3].IP)
	assert.Equal(t, "sameetandpotatoes", users[3].Username)
}

func DeleteTest(t *testing.T) {
	assert.True(t, DeleteUser("karateAMD"))
	assert.Equal(t, 5, len(QuerySessions()))
	assert.Equal(t, 5, len(QueryMessages()))
	assert.Equal(t, 6, len(QueryUsers()))
	assert.False(t, UserExists("karateAMD"))
}
