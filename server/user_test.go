package server

import "testing"
import (
	"github.com/stretchr/testify/assert"
)

func TestUser_Login(t *testing.T) {
	user := UserLogin("somelogin", "somepassword")
	assert.Nil(t, user)
}