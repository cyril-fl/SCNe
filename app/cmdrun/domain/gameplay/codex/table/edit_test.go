package codextable

import (
	"testing"
	"time"

	"scne/app/cmdrun/types/id"
)

func TestSimpleTablesEditEntryInPlace(t *testing.T) {
	locations := NewLocations()
	objects := NewObjects()
	quests := NewQuests()

	tests := []struct {
		name      string
		table     Table
		count     func() int
		timestamp func() time.Time
	}{
		{"locations", locations, func() int { return len(locations.Entries) }, func() time.Time { return locations.Entries[0].Timestamp }},
		{"objects", objects, func() int { return len(objects.Entries) }, func() time.Time { return objects.Entries[0].Timestamp }},
		{"quest", quests, func() int { return len(quests.Entries) }, func() time.Time { return quests.Entries[0].Timestamp }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.table.SaveFromMap(id.ID{}, map[string]string{"name": "entry", "description": "description"}); err != nil {
				t.Fatalf("adding an entry failed: %v", err)
			}
			entry := tt.table.FormEntries()[0]
			created := tt.timestamp()

			values := entry.Values
			values["name"] = "entry.edited"
			if err := tt.table.SaveFromMap(entry.ID, values); err != nil {
				t.Fatalf("editing the entry failed: %v", err)
			}

			edited := tt.table.FormEntries()[0]
			if tt.count() != 1 || edited.ID != entry.ID || edited.GetName() != "entry.edited" {
				t.Fatalf("after edit: count %d, entry %#v", tt.count(), edited)
			}
			if !tt.timestamp().Equal(created) {
				t.Fatalf("timestamp changed from %v to %v", created, tt.timestamp())
			}
		})
	}
}
