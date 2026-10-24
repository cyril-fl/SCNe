package codextable

import (
	"testing"

	"solopg/app/cmdrun/domain/card/characters/races"
)

func TestBeastiaryEditReplacesRaceInConfig(t *testing.T) {
	table := NewBeastiary()
	table.Entries = nil
	values := map[string]string{
		"name":        "beast.test",
		"description": "description",
		"encounter":   "encounter_normal",
	}
	if err := table.AddFromMap(values); err != nil {
		t.Fatalf("AddFromMap failed: %v", err)
	}
	entryID := table.Entries[0].ID

	unchanged := table.FormEntries()[0].Values
	unchanged["encounter"] = "encounter_normal"
	if err := table.SaveFromMap(entryID, unchanged); err != nil {
		t.Fatalf("saving an unchanged beast failed: %v", err)
	}

	edited := table.FormEntries()[0].Values
	edited["name"] = "beast.test.renamed"
	edited["encounter"] = "encounter_normal"
	if err := table.SaveFromMap(entryID, edited); err != nil {
		t.Fatalf("renaming a beast failed: %v", err)
	}

	if len(table.Entries) != 1 || table.Entries[0].ID != entryID {
		t.Fatalf("entries = %#v, want the edited entry only", table.Entries)
	}
	if races.FindByName("beast.test") != nil || races.FindByName("beast.test.renamed") == nil {
		t.Fatalf("races = %v, want beast.test replaced by beast.test.renamed", races.ListNames())
	}
}
