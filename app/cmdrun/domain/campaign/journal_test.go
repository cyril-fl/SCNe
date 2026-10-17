package campaign

import (
	"testing"
	"time"

	"solopg/app/cmdrun/types/id"
)

func TestJournalEnsureInitialized(t *testing.T) {
	journal := &Journal{}
	journal.EnsureInitialized()

	if journal.Entries == nil {
		t.Fatal("EnsureInitialized left Entries nil")
	}
}

func TestJournalAddEntry(t *testing.T) {
	journal := NewJournal(nil)
	entry := journal.AddEntry("gm", "An event")

	if len(journal.Entries) != 1 || journal.Entries[0] != *entry {
		t.Fatal("AddEntry did not append the entry")
	}
	if entry.ID == (id.ID{}) || entry.Author != "gm" || entry.Message != "An event" || entry.Timestamp.Location() != time.UTC {
		t.Fatal("AddEntry did not initialize the entry correctly")
	}
}

func TestEntryString(t *testing.T) {
	entry := Entry{Author: "gm", Message: "An event", Timestamp: time.Date(2026, time.January, 2, 15, 4, 0, 0, time.UTC)}

	if got, want := entry.String(), "Jan 02 15:04 - (gm) An event"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
