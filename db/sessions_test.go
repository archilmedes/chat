package db

import (
	"testing"
	"database/sql"
	"github.com/stretchr/testify/assert"
)

func SessionsTest (t *testing.T, db *sql.DB){
	sessions := QuerySessions(db)
	assert.Equal(t, 7, len(sessions))
}

