package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func SetupDatabaseForTests(t *testing.T) {
	SetupTestDatabase()
	tables := ShowTables()
	assert.Equal(t, 4, len(tables))
}
