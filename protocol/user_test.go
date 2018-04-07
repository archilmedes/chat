package protocol

import "testing"
import "github.com/stretchr/testify/assert"

func TestUser_IsSecure(t *testing.T) {
	secureUser := User{proto: OTRProtocol{}}
	assert.Equal(t, true, secureUser.IsSecure())
}

func TestNewUser(t *testing.T) {
	NewSecureUser()
}

func TestUser_Login(t *testing.T) {
	user := UserLogin("somelogin", "somepassword")
	assert.Nil(t, user)
}