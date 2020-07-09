package storage

import (
	"time"
)

type Source struct {
	SourceFile string
	Records    *[]Record
	SyncedAt   *time.Time
}

func (s *Source) Size() int {
	return len(*s.Records)
}
