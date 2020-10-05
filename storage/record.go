package storage

import (
	"time"
)

type Record struct {
	Date    time.Time
	Entries []Entry
}

// MostRecentEntry returns a pointer to the most recent entry
func (r *Record) MostRecentEntry() *Entry {
	if len(r.Entries) == 0 {
		return nil
	}

	lastIndex := len(r.Entries) - 1

	return &r.Entries[lastIndex]
}
