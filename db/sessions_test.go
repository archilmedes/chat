package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func SessionsTest(t *testing.T) {
	InitialSetupTest(t)
	InsertSessionTest(t)
	InsertServerDataTest(t)
	UpdateOtrDataTest(t)
	DeleteSessionTest(t)
}

func InitialSetupTest(t *testing.T) {
	sessions := QuerySessions()
	assert.Equal(t, 7, len(sessions))
	assert.Equal(t, "CBDEABD347ABDC392", sessions[1].Fingerprint)
	assert.Equal(t, "903873473785b4084787d8767889e988767543c45655434567a56467897669c009987766565b78765545e7896767f3234674734589bc9084734efb398473", sessions[2].PrivateKey)
	assert.Equal(t, 3, sessions[3].UserId)
	assert.Equal(t, 35, sessions[4].SSID)
	assert.Equal(t, 2, sessions[5].FriendId)
	assert.Equal(t, "B3847C837D77654E5", sessions[6].Fingerprint)
	assert.Equal(t, "bcd8763728749378ab8347839847328492ae897638903478b834743898c834738423e9786f9ff7657843cb3874383b8973487ef3864727384a8783647873", sessions[6].PrivateKey)
	assert.Equal(t, 6, sessions[6].UserId)
	assert.Equal(t, 64, sessions[6].SSID)
	assert.Equal(t, 4, sessions[6].FriendId)
}

func InsertSessionTest(t *testing.T) {
	assert.True(t, InsertIntoSessions(84, 8, 4, "abcdefabcdef1234567890abcdef1234567890123545734482", "ABC47364EDFB8664"))
	sessions := QuerySessions()
	assert.Equal(t, 8, len(sessions))
	assert.Equal(t, 84, sessions[7].SSID)
	assert.Equal(t, 4, sessions[7].FriendId)
	assert.Equal(t, "abcdefabcdef1234567890abcdef1234567890123545734482", sessions[7].PrivateKey)
}

func InsertServerDataTest(t *testing.T) {
	assert.True(t, InsertServerData(8, 1, 81))
	sessions := QuerySessions()
	assert.Equal(t, 9, len(sessions))
	assert.Equal(t, 8, sessions[7].UserId)
	assert.Equal(t, 1, sessions[7].FriendId)
	assert.Equal(t, "", sessions[7].Fingerprint)
}

func UpdateOtrDataTest(t *testing.T) {
	assert.True(t, UpdateSessionsOtrData(81, "16246374859ae8473743f343b314a76767d87677e8890d32223f", "ABD3474357DBE"))
	sessions := QuerySessions()
	assert.Equal(t, 9, len(sessions))
	assert.Equal(t, 8, sessions[7].UserId)
	assert.Equal(t, 1, sessions[7].FriendId)
	assert.Equal(t, "ABD3474357DBE", sessions[7].Fingerprint)
	assert.Equal(t, "16246374859ae8473743f343b314a76767d87677e8890d32223f", sessions[7].PrivateKey)
}

func DeleteSessionTest(t *testing.T) {
	assert.True(t, DeleteSession(81))
	sessions := QuerySessions()
	assert.Equal(t, 8, len(sessions))
}
