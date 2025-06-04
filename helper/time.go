package helper

import "time"

func ToISO(t time.Time) string {
	// Convert time to ISO 8601 format
	return t.Format(time.RFC3339)
}

func FromISO(isoStr string) time.Time {
	// Parse ISO 8601 formatted string to time.Time
	t, err := time.Parse(time.RFC3339, isoStr)
	if err != nil {
		return time.Time{}
	}
	return t
}
