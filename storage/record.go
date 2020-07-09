package storage

import (
	"time"
)

type Record struct {
	Date    time.Time
	Entries []Entry
}
