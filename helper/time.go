package helper

import "time"

func ToISO(t time.Time) string {
	// Convert time to ISO 8601 format
	return t.Format(time.RFC3339)
}
