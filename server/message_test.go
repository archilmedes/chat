package server

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const(
	sourceIP = "1.2.3.4"
	sourceMAC = "AA:BB:CC:DD:EE"
	sourceUsername = "Source Username"
	destUsername = "Dest Username"
)

func getGenericMessage() Message {
	msg := new(GenericMessage)
	msg.NewPayload(sourceMAC, sourceIP, sourceUsername, destUsername)
	return msg
}

func TestGenericMessage_DestID(t *testing.T) {
	msg := getGenericMessage()
	assert.Equal(t, destUsername, msg.DestID())
}

func TestGenericMessage_SourceID(t *testing.T) {
	msg := getGenericMessage()
	mac, ip, username := msg.SourceID()
	assert.Equal(t, sourceMAC, mac)
	assert.Equal(t, sourceIP, ip)
	assert.Equal(t, sourceUsername, username)
}