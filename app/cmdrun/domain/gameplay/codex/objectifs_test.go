package codex

import "testing"

func TestObjectivesFormEntriesAndUpdate(t *testing.T) {
	table := NewObjectivesTable([]ObjectifEntry{{Title: "old", Description: "description"}})
	entries := table.FormEntries()
	if len(entries) != 1 || entries[0].Values["title"] != "old" {
		t.Fatalf("unexpected form entry: %#v", entries)
	}

	if err := table.SaveFromMappedValues(0, map[string]string{
		"title": "new", "description": "updated",
	}); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if table.Entries[0].Title != "new" || table.Entries[0].Description != "updated" {
		t.Fatalf("entry was not updated: %#v", table.Entries[0])
	}
}
