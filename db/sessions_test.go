package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func SessionsTest(t *testing.T) {
	sessions := QuerySessions()
	assert.Equal(t, 7, len(sessions))
	assert.Equal(t, "CBDEABD347ABDC392", sessions[1].fingerprint)
	assert.Equal(t, "903873473785b4084787d8767889e988767543c45655434567a56467897669c009987766565b78765545e7896767f3234674734589bc9084734efb398473", sessions[2].privateKey)
	assert.Equal(t, 3, sessions[3].userId)
	assert.Equal(t, 35, sessions[4].SSID)
	assert.Equal(t, 2, sessions[5].friendId)
	assert.Equal(t, "B3847C837D77654E5", sessions[6].fingerprint)
	assert.Equal(t, "bcd8763728749378ab8347839847328492ae897638903478b834743898c834738423e9786f9ff7657843cb3874383b8973487ef3864727384a8783647873", sessions[6].privateKey)
	assert.Equal(t, 6, sessions[6].userId)
	assert.Equal(t, 64, sessions[6].SSID)
	assert.Equal(t, 4, sessions[6].friendId)
}
