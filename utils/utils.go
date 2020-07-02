package utils

import "time"

// Today returns midnight for the current date
func Today() *time.Time {
	y1, m1, d1 := time.Now().Date()

	today := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	return &today
}
