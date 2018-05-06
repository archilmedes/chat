package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetFormattedTime(t *testing.T) {
	someTime := time.Unix(100, 1000)
	assert.Equal(t, "1969-12-31 18:01:40", GetFormattedTime(someTime))
}
