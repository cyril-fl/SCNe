package codextable

import (
	"errors"
	"fmt"
	"os"
	"scne/app/cmdrun/domain/card/attributes/rarity"
	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/attributes/variety"
	"scne/app/cmdrun/domain/card/locations"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/logs"
	"slices"
	"strconv"
	"time"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// - Table --//

// LocationsTable lists the locations of the campaign.
type LocationsTable struct {
	*Data[locationsentry] `bson:",inline"`
}

type locationsentry struct {
	ID        id.ID               `bson:"id"`
	Timestamp time.Time           `bson:"timestamp"`
	Location  *locations.Location `bson:"location"`
}

var tablelocations = "codex.locations"

// NewLocations returns a location table with entries.
func NewLocations(entries ...locationsentry) *LocationsTable {
	return &LocationsTable{
		Data: new(i19n.Localize(tablelocations), entries),
	}
}

// Getters & Setters

// Methods

// AddFromMap creates a location from form values and appends it.
func (tb *LocationsTable) AddFromMap(values map[string]string) error {
	location, err := tb.locationFromMap(values)
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

func (tb *LocationsTable) locationFromMap(values map[string]string) (*locations.Location, error) {
	tb.assertEntry(values)

	if tb.HasErr() {
		cwd, cwderr := os.Getwd()
		return nil, logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, tb.GetErr()),
		})
	}

	location, err := locations.New(locations.Template{
		Name:        values["name"],
		Description: values["description"],
		Rarity:      rarity.Default(),
		Variety:     variety.MakeDefault("location_card"),
		Effects:     effectsFromValues(values),
	})

	if err != nil {
		return nil, err
	}

	return location, nil
}

// effectsFromValues reads the effect.<i>.* form groups: each filled stat becomes one effect.
func effectsFromValues(values map[string]string) []stats.Effect {
	effects := []stats.Effect{}
	for i := 0; ; i++ {
		prefix := fmt.Sprintf("effect.%d.", i)
		name, ok := values[prefix+"name"]
		if !ok {
			return effects
		}

		for _, stat := range stats.List() {
			raw := values[prefix+stat.String()]
			if raw == "" {
				continue
			}

			value, _ := strconv.Atoi(raw)
			effects = append(effects, stats.Effect{
				Descriptible: plumprimitive.Describe(plumprimitive.TemplateDescriptible{
					Name:        name,
					Description: values[prefix+"description"],
				}),
				Modifier: stats.Modifier{Stat: stat, Value: value},
			})
		}
	}
}

// Add appends a new location entry.
func (tb *LocationsTable) Add(location *locations.Location) {
	tb.Entries = append(tb.Entries, locationsentry{
		ID:        id.New(),
		Timestamp: time.Now().UTC(),
		Location:  location,
	})
}

// FormEntries returns the entries as form page rows.
func (tb *LocationsTable) FormEntries() []Entry {
	entries := make([]Entry, 0, len(tb.Entries))
	for _, entry := range tb.Entries {
		if entry.Location == nil {
			continue
		}
		name, description := entry.Location.GetName(), entry.Location.GetDescription()
		values := map[string]string{
			"name":        name,
			"description": description,
		}
		for i, effect := range entry.Location.Effects {
			prefix := fmt.Sprintf("effect.%d.", i)
			values[prefix+"name"] = effect.GetName()
			values[prefix+"description"] = effect.GetDescription()
			values[prefix+effect.Modifier.Stat.String()] = strconv.Itoa(effect.Modifier.Value)
		}
		entries = append(entries, newEntry(entry.ID, name, description, values, locations.IsLocked(name)))
	}
	return entries
}

// Find returns the location of the entry with entryID, or nil.
func (tb *LocationsTable) Find(entryID id.ID) *locations.Location {
	index := slices.IndexFunc(tb.Entries, func(entry locationsentry) bool { return entry.ID == entryID })
	if index < 0 {
		return nil
	}
	return tb.Entries[index].Location
}

// SaveFromMap edits the entry with entryID, or adds one when it is zero or unknown.
func (tb *LocationsTable) SaveFromMap(entryID id.ID, values map[string]string) error {
	index := slices.IndexFunc(tb.Entries, func(entry locationsentry) bool { return entry.ID == entryID })
	if entryID.IsZero() || index < 0 {
		return tb.AddFromMap(values)
	}
	return tb.editFromMap(index, values)
}

func (tb *LocationsTable) editFromMap(index int, values map[string]string) error {
	location, err := tb.locationFromMap(values)
	if err != nil {
		return err
	}

	tb.Entries[index].Location = location
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

// Ensure initializes the table when missing.
func (tb *LocationsTable) Ensure() {
	ensureEmbeddedTable(i19n.Localize(tablelocations), &tb.Data)
}

// - Helper --//
