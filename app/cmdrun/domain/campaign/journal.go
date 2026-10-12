package campaign

import (
	"solopg/app/cmdrun/types/id"
	"time"
)

type Journal struct {
	Entries []Entry
}

type Entry struct {
	ID        id.ID
	Author    string
	Message   string
	Timestamp time.Time
}

func NewJournal(entries []Entry) *Journal {
	return &Journal{
		Entries: entries,
	}
}

func (j *Journal) EnsureInitialized() {
	if j == nil {
		j = NewJournal([]Entry{})
	}

	if j.Entries == nil {
		j.Entries = []Entry{}
	}
}

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

// String format for the entry in the journal, shown in gamebord view
func (e *Entry) String() string {
	return e.Timestamp.Format("Jan 02 15:04") + " - (" + e.Author + ") " + e.Message
}
