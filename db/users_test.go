package db

import (
	"testing"
	"database/sql"
	"github.com/stretchr/testify/assert"
)

func UsersTest (t *testing.T, db *sql.DB){
	users := QueryUsers(db)
	assert.Equal(t, 5, len(users))
	assert.Equal(t, "alicepassword", users[0].password)
	assert.Equal(t, "bob", users[1].login)
	assert.Equal(t, 3, users[2].id)
	assert.Equal(t, 4, users[3].id)
	assert.Equal(t, "sameet", users[3].login)
	assert.Equal(t, "iLuvMacs", users[3].password)


}

