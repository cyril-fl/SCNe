package campaign

import (
	"testing"
	"time"

	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/i19n"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
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

func setTranslator(t *testing.T, messages ...*i18n.Message) {
	t.Helper()
	bundle := i18n.NewBundle(language.English)
	bundle.AddMessages(language.English, messages...)
	translator := i19n.New()
	translator.SetLocale(i18n.NewLocalizer(bundle, "en"))
	i19n.SetCache(translator)
	t.Cleanup(func() { i19n.SetCache(nil) })
}

func TestEntryString(t *testing.T) {
	setTranslator(t, &i18n.Message{ID: "journal.date:format", Other: "02/01 15:04"})

	entry := Entry{Author: "gm", Message: "An event", Timestamp: time.Date(2026, time.January, 2, 15, 4, 0, 0, time.UTC)}

	if got, want := entry.String(), "02/01 15:04 - (gm) An event"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
	entry.Author = ""
	if got, want := entry.String(), "02/01 15:04 - An event"; got != want {
		t.Fatalf("String() without author = %q, want %q", got, want)
	}
}

func TestDisplayAuthor(t *testing.T) {
	setTranslator(t, &i18n.Message{ID: "author.system", Other: "Système"})

	if got, want := DisplayAuthor(SYSTEM_AUTHOR), "Système"; got != want {
		t.Fatalf("DisplayAuthor(system) = %q, want %q", got, want)
	}
	if got, want := DisplayAuthor("Ada"), "Ada"; got != want {
		t.Fatalf("DisplayAuthor(player) = %q, want %q", got, want)
	}
}
