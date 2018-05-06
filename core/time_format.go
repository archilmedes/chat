package core

import (
	"time"
)

// Gets the formatted input time to save in the database
func GetFormattedTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
