package codextable

import (
	"testing"

	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/card/characters/classes"
	"scne/app/cmdrun/domain/card/characters/races"
)

func newTestNpcs(t *testing.T, class, race string) (*NpcsTable, map[string]string) {
	t.Helper()
	initialStats := make(stats.Stats)
	for _, stat := range stats.List() {
		initialStats[stat] = 10
	}
	character, err := characters.New(characters.Template{
		Name:        "npc.test",
		Description: "description",
		Class:       class,
		Race:        race,
		Stats:       initialStats,
	})
	if err != nil {
		t.Fatalf("characters.New failed: %v", err)
	}

	table := NewNpcs()
	table.Add(character)
	return table, table.FormEntries()[0].Values
}

func TestNpcsEditDoesNotStackStats(t *testing.T) {
	classList, raceList := classes.List(), races.List()
	if len(stats.List()) == 0 || len(classList) == 0 || len(raceList) == 0 {
		t.Skip("no stat, class or race in the loaded config")
	}
	table, values := newTestNpcs(t, classList[0].GetName(), raceList[0].GetID())
	entryID := table.Entries[0].ID

	for range 2 {
		if err := table.SaveFromMap(entryID, values); err != nil {
			t.Fatalf("SaveFromMap failed: %v", err)
		}
	}

	if len(table.Entries) != 1 || table.Entries[0].ID != entryID {
		t.Fatalf("entries = %#v, want the edited entry only", table.Entries)
	}
	for stat, value := range table.Entries[0].Character.Stats {
		if value != 10 {
			t.Fatalf("stat %s = %d after edits, want 10", stat, value)
		}
	}
}

func TestNpcsEditSwapsClassBonus(t *testing.T) {
	classList, raceList := classes.List(), races.List()
	if len(stats.List()) == 0 || len(classList) < 2 || len(raceList) == 0 {
		t.Skip("needs stats, two classes and a race in the loaded config")
	}
	previousClass, editedClass := classList[0], classList[1]
	table, values := newTestNpcs(t, previousClass.GetName(), raceList[0].GetID())
	values["class"] = editedClass.GetName()

	if err := table.SaveFromMap(table.Entries[0].ID, values); err != nil {
		t.Fatalf("SaveFromMap failed: %v", err)
	}

	want := make(stats.Stats)
	for _, stat := range stats.List() {
		want[stat] = 10
	}
	for _, modifier := range previousClass.GetBonus() {
		want[modifier.Stat] -= modifier.Value
	}
	for _, modifier := range editedClass.GetBonus() {
		want[modifier.Stat] += modifier.Value
	}
	edited := table.Entries[0].Character
	if edited.Class != editedClass.GetName() {
		t.Fatalf("class = %q, want %q", edited.Class, editedClass.GetName())
	}
	for stat, value := range want {
		if edited.Stats[stat] != value {
			t.Fatalf("stat %s = %d, want %d", stat, edited.Stats[stat], value)
		}
	}
}
