package db

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func UsersTest (t *testing.T){
	users := QueryUsers()
	UserSetupTest(t, users)
	UserExistsTest(t)
	UserExistsAfterDeleteTest(t)
	ValidateCredentialsTest(t)

}

func ValidateCredentialsTest(t *testing.T) {
	assert.True(t, ValidateCredentials("karateAMD", "pwd123"))
	assert.False(t, ValidateCredentials("Sameet", "iLuvMacs"))
	assert.False(t, ValidateCredentials("sameetandpotatoes", "linuxFTW"))
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
	assert.Equal(t, 5, len(users))
	assert.Equal(t, "alicepassword", users[0].password)
	assert.Equal(t, "bob", users[1].username)
	assert.Equal(t, "192.168.10.123", users[2].ipAddress)
	assert.Equal(t, "10.192.345.987", users[3].ipAddress)
	assert.Equal(t, "sameetandpotatoes", users[3].username)
	assert.Equal(t, "iLuvMacs", users[3].password)
}

