package codex

import (
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/cmdrun/domain/card/locations"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/utils/transform"
	"time"
)

// - Table --//
type LocationsTable struct {
	*TableData[LocationsEntry]
}

type LocationsEntry struct {
	Timestamp time.Time           `bson:"timestamp"`
	Location  *locations.Location `bson:"location"`
}

var tablelocations = "codex.locations"

func NewLocationsTable(entries []LocationsEntry) *LocationsTable {
	return &LocationsTable{
		TableData: newTable(i19n.Localize(tablelocations), entries),
	}
}

// - Methods --//
func (tb *LocationsTable) Add(location *locations.Location) {
	tb.Entries = append(tb.Entries, LocationsEntry{
		Timestamp: time.Now().UTC(),
		Location:  location,
	})
}

func (tb *LocationsTable) AddFromMappedValues(values map[string]string) error {
	if err := tb.assertEntry(values); err != nil {
		return err
	}

	location, err := locations.New(locations.Template{
		Name:        values["name"],
		Description: values["description"],
		Rarity:      rarity.Default(),
		Variety:     variety.MakeDefault("location_card"),
	})

	if err != nil {
		return err
	}

	tb.Add(location)

	logs.SilentSuccess("codex.entry:added", map[string]any{
		"Table": tablelocations,
		"Entry": transform.ParseJson(location),
	})

	return nil
}

func (tb *LocationsTable) FormEntries() []Entry {
	entries := make([]Entry, 0, len(tb.Entries))
	for _, entry := range tb.Entries {
		if entry.Location == nil {
			continue
		}
		name, description := entry.Location.GetName(), entry.Location.GetDescription()
		entries = append(entries, newEntry(name, description, map[string]string{
			"name": name, "description": description,
		}, locations.IsLocked(name)))
	}
	return entries
}

func (tb *LocationsTable) SaveFromMappedValues(index int, values map[string]string) error {
	if index < 0 {
		return tb.AddFromMappedValues(values)
	}
	oldLen := len(tb.Entries)
	if err := tb.AddFromMappedValues(values); err != nil {
		return err
	}
	if index < oldLen {
		tb.Entries[index] = tb.Entries[oldLen]
		tb.Entries = tb.Entries[:oldLen]
	}
	return nil
}

func (tb *LocationsTable) FindEntryByName(name string) *LocationsEntry {
	for i, entry := range tb.Entries {
		if entry.Location != nil && entry.Location.GetName() == name {
			return &tb.Entries[i]
		}
	}
	return nil
}

// - Helper --//

func (tb *LocationsTable) assertEntry(values map[string]string) error {
	var err []error

	if values["name"] == "" {

		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "location",
			"Property": "name",
		}))
	}
	if values["description"] == "" {
		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "location",
			"Property": "description",
		}))
	}

	if len(err) > 0 {
		return tb.formatAssertErrors(err)
	}

	return nil
}

func (tb *LocationsTable) Ensure() {
	if tb.name == "" {
		tb.name = i19n.Localize(tablelocations)
	}
	if tb.Entries == nil {
		tb.Entries = []LocationsEntry{}
	}
}
