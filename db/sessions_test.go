package db

import (
	"testing"
	"database/sql"
	"github.com/stretchr/testify/assert"
)

func SessionsTest (t *testing.T, db *sql.DB){
	sessions := QuerySessions(db)
	assert.Equal(t, 7, len(sessions))
	assert.Equal(t, "CBDEABD347ABDC392", sessions[1].fingerprint)
	assert.Equal(t, "bcd8763728749378ab8347839847328492ae897638903478b834743898c834738423e9786f9ff7657843cb3874383b8973487ef3864727384a8783647873", sessions[2].privateKey)
	assert.Equal(t, 3, sessions[3].userId)
	assert.Equal(t, 35, sessions[4].SSID)
	assert.Equal(t, 2, sessions[5].friendId)
	assert.Equal(t, "675C6A7CA877B6A67", sessions[6].fingerprint)
	assert.Equal(t, "abc78384689234752369071625d8736543976d7f967798b567789098768789076890876890786890876890e8789087e877a90887b89c87d77e762532dbcf", sessions[6].privateKey)
	assert.Equal(t, 5, sessions[6].userId)
	assert.Equal(t, 52, sessions[6].SSID)
	assert.Equal(t, 2, sessions[6].friendId)
}

