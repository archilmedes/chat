package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetAddresses(t *testing.T) {
	mac, ip, err := GetAddresses()
	assert.Nil(t, err)
	assert.NotNil(t, ip)
	assert.NotNil(t, mac)
}