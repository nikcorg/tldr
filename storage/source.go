package storage

import (
	"time"
)

type Source struct {
	SourceFile string
	Records    []*Record
	SyncedAt   *time.Time
}

func (s *Source) Size() int {
	return len(s.Records)
}

func (s *Source) WasSynced() bool {
	return s.SyncedAt != nil
}

func (s *Source) FirstRecord() *Record {
	return s.Records[0]
}
