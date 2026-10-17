package codextable

import (
	"errors"
	"testing"
)

func TestNewEntryAndFindEntryByName(t *testing.T) {
	entry := newEntry("name", "description", map[string]string{"name": "name"}, false)

	found := findEntryByName([]Entry{entry}, "name")

	if found == nil {
		t.Fatal("findEntryByName did not find the entry")
	}
	if found.GetName() != "name" {
		t.Fatalf("found entry has name %q, want %q", found.GetName(), "name")
	}
	if findEntryByName([]Entry{entry}, "missing") != nil {
		t.Fatal("findEntryByName found a missing entry")
	}
}

func TestEnsureEmbeddedTableInitializesOnlyMissingData(t *testing.T) {
	var table *data[string]

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
		{name: "locations", table: NewLocationsTable(nil)},
		{name: "npcs", table: NewNpcsTable(nil)},
		{name: "objects", table: NewObjectsTable(nil)},
		{name: "quest", table: NewQuestTable(nil)},
		{name: "beastiary", table: NewBeastiaryTable(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.table == nil {
				t.Fatal("constructor returned nil")
			}
			if entries := tt.table.FormEntries(); entries == nil {
				t.Fatal("FormEntries returned nil for an empty table")
			}
			if tt.table.FindEntryByName("missing") != nil {
				t.Fatal("FindEntryByName found an entry in an empty table")
			}
		})
	}
}
