package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func UsersTest(t *testing.T) {
	users := QueryUsers()
	UserSetupTest(t, users)
	UserExistsTest(t)
	UserExistsAfterDeleteTest(t)
	ValidateCredentialsTest(t)
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

func UserExistsAfterDeleteTest(t *testing.T) {
	assert.True(t, DeleteUser("sameetandpotatoes"))
	assert.False(t, UserExists("sameetandpotatoes"))
}

func UserSetupTest(t *testing.T, users []User) {
	assert.Equal(t, 6, len(users))
	assert.Equal(t, "bob", users[1].Username)
	assert.Equal(t, "192.168.10.123", users[2].IP)
	assert.Equal(t, "10.192.345.987", users[3].IP)
	assert.Equal(t, "sameetandpotatoes", users[3].Username)
}
