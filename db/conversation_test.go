package db

import (
	"testing"
	"database/sql"
	"github.com/stretchr/testify/assert"
)

func ConversationTest (t *testing.T, db *sql.DB){
	conversations := QueryConversations(db)
	assert.Equal(t, 8, len(conversations))
}

