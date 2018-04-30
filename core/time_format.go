package core

import (
	"fmt"
	"strings"
	"time"
)

// Gets the formatted input time to save in the database
func GetFormattedTime(t time.Time) string {
	timestampParts := strings.Split(t.String(), " ")
	return fmt.Sprintf("%s %s", timestampParts[0], timestampParts[1])
}
