package codextable

import (
	"errors"
	"os"
	"testing"

	"solopg/app/cmdrun/types/id"
	"solopg/app/shared/services/logs"
)

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

func TestMain(m *testing.M) {
	logs.Init(testLogRegistrar{})
	code := m.Run()
	logs.Init(nil)
	os.Exit(code)
}

func TestNewEntryAndFindEntryByName(t *testing.T) {
	entry := newEntry(id.New(), "name", "description", map[string]string{"name": "name"}, false)

	found := FindEntryByName([]Entry{entry}, "name")

	if found == nil {
		t.Fatal("findEntryByName did not find the entry")
	}
	if found.GetName() != "name" {
		t.Fatalf("found entry has name %q, want %q", found.GetName(), "name")
	}
	if FindEntryByName([]Entry{entry}, "missing") != nil {
		t.Fatal("FindEntryByName found a missing entry")
	}
}

func TestEnsureEmbeddedTableInitializesOnlyMissingData(t *testing.T) {
	var table *Data[string]

	ensureEmbeddedTable("table", &table)

	if table == nil {
		t.Fatal("ensureEmbeddedTable did not initialize the table")
	}
	if table.name != "table" {
		t.Fatalf("table name is %q, want %q", table.name, "table")
	}
	if table.Entries == nil {
		t.Fatal("ensureEmbeddedTable initialized nil entries")
	}

	existing := new("existing", []string{"entry"})
	ensureEmbeddedTable("replacement", &existing)

	if existing.name != "existing" {
		t.Fatalf("ensureEmbeddedTable replaced existing data with %q", existing.name)
	}
	if len(existing.Entries) != 1 || existing.Entries[0] != "entry" {
		t.Fatal("ensureEmbeddedTable changed existing entries")
	}
}

func TestTableDataEmbedsFallible(t *testing.T) {
	table := new[string]("table", nil)
	table.SetErr(errors.New("invalid"))

	if !table.HasErr() || table.GetErr() == nil {
		t.Fatal("table data does not expose Fallible behavior")
	}

	table.ClearErr()
	if table.HasErr() || table.GetErr() != nil {
		t.Fatal("table data did not clear its error state")
	}
}

func TestSubtablesExposeEmptyFormEntries(t *testing.T) {
	tests := []struct {
		name  string
		table Table
	}{
		{name: "locations", table: NewLocations()},
		{name: "npcs", table: NewNpcs()},
		{name: "objects", table: NewObjects()},
		{name: "quest", table: NewQuests()},
		{name: "beastiary", table: NewBeastiary()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.table == nil {
				t.Fatal("constructor returned nil")
			}
			entries := tt.table.FormEntries()
			if entries == nil {
				t.Fatal("FormEntries returned nil for an empty table")
			}
			if FindEntryByName(entries, "missing") != nil {
				t.Fatal("FindEntryByName found an entry in an empty table")
			}
		})
	}
}

func TestSaveFromMapReplacesEntryByID(t *testing.T) {
	table := NewQuests()
	for _, name := range []string{"first", "second"} {
		if err := table.AddFromMap(map[string]string{"name": name, "description": "description"}); err != nil {
			t.Fatalf("AddFromMap(%q) error = %v", name, err)
		}
	}
	// A nil entry is skipped by FormEntries, so form positions no longer match table positions.
	table.Entries = append([]questentry{{ID: id.New()}}, table.Entries...)
	target := table.Entries[2].ID

	err := table.SaveFromMap(target, map[string]string{"name": "edited", "description": "description"})
	if err != nil {
		t.Fatalf("SaveFromMap() error = %v", err)
	}

	if len(table.Entries) != 3 {
		t.Fatalf("table has %d entries, want 3", len(table.Entries))
	}
	if table.Entries[1].Quest.GetName() != "first" {
		t.Fatalf("untouched entry renamed to %q", table.Entries[1].Quest.GetName())
	}
	if table.Entries[2].Quest.GetName() != "edited" || table.Entries[2].ID != target {
		t.Fatalf("edited entry = %q %v, want %q %v", table.Entries[2].Quest.GetName(), table.Entries[2].ID, "edited", target)
	}
}

func TestSaveFromMapAddsEntryForZeroID(t *testing.T) {
	table := NewQuests()

	if err := table.SaveFromMap(id.ID{}, map[string]string{"name": "new", "description": "description"}); err != nil {
		t.Fatalf("SaveFromMap() error = %v", err)
	}

	if len(table.Entries) != 1 || table.Entries[0].ID.IsZero() {
		t.Fatalf("entries = %#v, want one entry with an ID", table.Entries)
	}
}

func TestEnsureAssignsMissingEntryIDs(t *testing.T) {
	table := NewQuests(questentry{})

	table.Ensure()

	if table.Entries[0].ID.IsZero() {
		t.Fatal("Ensure did not assign an ID to a legacy entry")
	}
}
