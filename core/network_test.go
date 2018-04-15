package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAddresses(t *testing.T) {
	mac, ip, err := GetAddresses()
	assert.Nil(t, err)
	assert.NotNil(t, ip)
	assert.NotNil(t, mac)
}
