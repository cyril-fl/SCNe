package campaign

import (
	"solopg/app/cmdrun/types/id"
	"time"
)

// Journal is an ordered list of authored entries.
type Journal struct {
	Entries []Entry
}

// Entry is a journal message.
type Entry struct {
	ID        id.ID
	Author    string
	Message   string
	Timestamp time.Time
}

// NewJournal returns a Journal with entries.
func NewJournal(entries []Entry) *Journal {
	return &Journal{
		Entries: entries,
	}
}

// EnsureInitialized replaces nil Entries with an empty slice; it does nothing on a nil j.
func (j *Journal) EnsureInitialized() {
	if j == nil {
		j = NewJournal([]Entry{})
	}

	if j.Entries == nil {
		j.Entries = []Entry{}
	}
}

// AddEntry appends a UTC-timestamped entry and returns it.
func (j *Journal) AddEntry(author, message string) *Entry {
	entry := &Entry{
		ID:        id.New(),
		Author:    author,
		Message:   message,
		Timestamp: time.Now().UTC(),
	}

	j.Entries = append(j.Entries, *entry)
	return entry
}

// String formats the entry as a gameboard journal line.
func (e *Entry) String() string {
	return e.Timestamp.Format("Jan 02 15:04") + " - (" + e.Author + ") " + e.Message
}
