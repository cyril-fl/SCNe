package campaign

import (
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/i19n"
	"time"
)

// SYSTEM_AUTHOR is the author saved for system entries, translated only on display.
const SYSTEM_AUTHOR = "System"

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

// RemoveLast drops the last entry; it does nothing on an empty journal.
func (j *Journal) RemoveLast() {
	if len(j.Entries) > 0 {
		j.Entries = j.Entries[:len(j.Entries)-1]
	}
}

// String formats the entry as a gameboard journal line, without author when it has none.
func (e *Entry) String() string {
	date := e.Timestamp.Format(i19n.Localize("journal.date:format")) + " - "
	if e.Author == "" {
		return date + e.Message
	}
	return date + "(" + DisplayAuthor(e.Author) + ") " + e.Message
}

// DisplayAuthor returns the name to show for author.
func DisplayAuthor(author string) string {
	if author == SYSTEM_AUTHOR {
		return i19n.Localize("author.system")
	}
	return author
}
