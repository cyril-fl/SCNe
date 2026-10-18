package codextable

import (
	"errors"
	"os"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/cmdrun/domain/card/locations"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"time"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// - Table --//
type LocationsTable struct {
	*data[locationsentry]
}

type locationsentry struct {
	Timestamp time.Time           `bson:"timestamp"`
	Location  *locations.Location `bson:"location"`
}

var tablelocations = "codex.locations"

func NewLocations(entries ...locationsentry) *LocationsTable {
	return &LocationsTable{
		data: new(i19n.Localize(tablelocations), entries),
	}
}

// Getters & Setters

// Methods
func (tb *LocationsTable) AddFromMap(values map[string]string) error {
	tb.assertEntry(values)

	if tb.HasErr() {
		cwd, cwderr := os.Getwd()
		return logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, tb.GetErr()),
		})
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
		"Entry": plumtransform.ParseJson(location),
	})

	return nil
}

func (tb *LocationsTable) Add(location *locations.Location) {
	tb.Entries = append(tb.Entries, locationsentry{
		Timestamp: time.Now().UTC(),
		Location:  location,
	})
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

func (tb *LocationsTable) SaveFromMap(index int, values map[string]string) error {
	if index < 0 {
		return tb.AddFromMap(values)
	}
	oldLen := len(tb.Entries)
	if err := tb.AddFromMap(values); err != nil {
		return err
	}
	if index < oldLen {
		tb.Entries[index] = tb.Entries[oldLen]
		tb.Entries = tb.Entries[:oldLen]
	}
	return nil
}

func (tb *LocationsTable) assertEntry(values map[string]string) {
	tb.ClearErr()

	if values["name"] == "" {
		tb.SetErr(logs.Error("error.required", map[string]any{
			"Subject":  "location",
			"Property": "name",
		}))
	}

	if values["description"] == "" {
		tb.SetErr(logs.Error("error.required", map[string]any{
			"Subject":  "location",
			"Property": "description",
		}))
	}
}

func (tb *LocationsTable) Ensure() {
	ensureEmbeddedTable(i19n.Localize(tablelocations), &tb.data)
}

// - Helper --//
