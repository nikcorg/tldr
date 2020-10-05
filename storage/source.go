package storage

import (
	"time"
)

// Source represents a set of records on disk
type Source struct {
	SourceFile string
	Records    []*Record
	SyncedAt   *time.Time
}

// Size returns the number of records it contains
func (s *Source) Size() int {
	return len(s.Records)
}

// WasSynced indicates whether a source has been synced
func (s *Source) WasSynced() bool {
	return s.SyncedAt != nil
}

// FirstRecord returns the first record
func (s *Source) FirstRecord() *Record {
	return s.Records[0]
}
