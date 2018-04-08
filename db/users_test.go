package db

import (
	"testing"
	"database/sql"
	"github.com/stretchr/testify/assert"
)

func UsersTest (t *testing.T, db *sql.DB){
	users := QueryUsers(db)
	assert.Equal(t, 5, len(users))
}

