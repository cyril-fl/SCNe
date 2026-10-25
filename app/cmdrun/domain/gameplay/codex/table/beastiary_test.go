package codextable

import (
	"strings"
	"testing"

	"scne/app/cmdrun/domain/card/characters/classes"
	"scne/app/cmdrun/domain/card/characters/races"
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
	raceID := table.Entries[0].Race.GetID()
	if got := races.NameOf(raceID); got != "beast.test.renamed" {
		t.Fatalf("race %s name = %q, want beast.test.renamed", raceID, got)
	}
}

func TestBeastiaryRenameIsFollowedByNpcs(t *testing.T) {
	if len(classes.List()) == 0 {
		t.Skip("no class in the loaded config")
	}
	beasts := NewBeastiary()
	beasts.Entries = nil
	if err := beasts.AddFromMap(map[string]string{
		"name":        "beast.followed",
		"description": "description",
		"encounter":   "encounter_normal",
	}); err != nil {
		t.Fatalf("AddFromMap failed: %v", err)
	}
	raceID := beasts.Entries[0].Race.GetID()
	npcs, values := newTestNpcs(t, classes.List()[0].GetName(), raceID)

	renamed := beasts.FormEntries()[0].Values
	renamed["name"] = "beast.followed.renamed"
	renamed["encounter"] = "encounter_normal"
	if err := beasts.SaveFromMap(beasts.Entries[0].ID, renamed); err != nil {
		t.Fatalf("renaming a beast failed: %v", err)
	}

	if got := npcs.FormEntries()[0].Description; !strings.Contains(got, "beast.followed.renamed") {
		t.Fatalf("npc race shown as %q, want beast.followed.renamed", got)
	}
	if err := npcs.SaveFromMap(npcs.Entries[0].ID, values); err != nil {
		t.Fatalf("editing the npc after the rename failed: %v", err)
	}
	for stat, value := range npcs.Entries[0].Character.Stats {
		if value != 10 {
			t.Fatalf("stat %s = %d after edit, want 10", stat, value)
		}
	}
}
