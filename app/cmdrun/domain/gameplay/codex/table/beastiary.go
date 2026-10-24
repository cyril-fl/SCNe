package codextable

import (
	"errors"
	"os"
	"slices"
	"solopg/app/cmdrun/domain/card/characters/races"
	"solopg/app/cmdrun/services/factory/millstats"
	"solopg/app/cmdrun/types/id"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"time"

	plumtypes "github.com/cyril-fl/plum/types"
	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// - Table --//
type BeastiaryTable struct {
	*Data[beastentry] `bson:",inline"`
}

type beastentry struct {
	ID        id.ID       `bson:"id"`
	Timestamp time.Time   `bson:"timestamp"`
	Race      *races.Race `bson:"race"`
}

var tablebeast = "codex.beastiary"

func NewBeastiary(entries ...beastentry) *BeastiaryTable {
	entries = superChargeEntries(entries)

	return &BeastiaryTable{
		Data: new(i19n.Localize(tablebeast), entries),
	}
}

// Getters & Setters

// Methods
func (tb *BeastiaryTable) Add(race *races.Race) {
	tb.Entries = append(tb.Entries, beastentry{
		ID:        id.New(),
		Timestamp: time.Now().UTC(),
		Race:      race,
	})
}

func (tb *BeastiaryTable) AddFromMap(values map[string]string) error {
	raw, err := tb.templateFromMap(values)
	if err != nil {
		return err
	}

	if err := races.AddInConfig(raw); err != nil {
		return logs.Error("error.unexpected", map[string]any{
			"Subject": "race",
			"Error":   err,
		})
	}

	race := races.New(raw)

	tb.Add(&race)

	logs.SilentSuccess("codex.entry:added", map[string]any{
		"Table": tablebeast,
		"Entry": plumtransform.ParseJson(race),
	})

	return nil
}

func (tb *BeastiaryTable) FormEntries() []Entry {
	entries := make([]Entry, 0, len(tb.Entries))
	for _, entry := range tb.Entries {
		if entry.Race == nil {
			continue
		}
		name, description := entry.Race.GetName(), entry.Race.GetDescription()
		values := map[string]string{"name": name, "description": description}
		for _, modifier := range entry.Race.GetBonus() {
			values[modifier.Stat.String()] = plumtransform.ParseAsString(modifier.Value)
		}
		entries = append(entries, newEntry(entry.ID, name, description, values, races.IsLocked(name)))
	}
	return entries
}

func (tb *BeastiaryTable) SaveFromMap(entryID id.ID, values map[string]string) error {
	index := slices.IndexFunc(tb.Entries, func(entry beastentry) bool { return entry.ID == entryID })
	if entryID.IsZero() || index < 0 {
		return tb.AddFromMap(values)
	}
	return tb.editFromMap(index, values)
}

// The edited beast replaces its race in the config instead of adding a new one.
// TODO renaming a beast leaves PNJs with the old race name: their next edit
// fails race validation and the old race bonus is never subtracted.
// Either rename the race on PNJs too, or make the name read-only on edit.
func (tb *BeastiaryTable) editFromMap(index int, values map[string]string) error {
	raw, err := tb.templateFromMap(values)
	if err != nil {
		return err
	}

	previous := tb.Entries[index].Race
	if err := races.ReplaceInConfig(previous.GetName(), raw); err != nil {
		return logs.Error("error.unexpected", map[string]any{
			"Subject": "race",
			"Error":   err,
		})
	}

	race := races.New(raw)
	tb.Entries[index].Race = &race
	return nil
}

func (tb *BeastiaryTable) templateFromMap(values map[string]string) (races.Template, error) {
	tb.assertEntry(values)

	if tb.HasErr() {
		cwd, cwderr := os.Getwd()
		return races.Template{}, logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, tb.GetErr()),
		})
	}

	generator := millstats.NewWithValues(values)
	generator.GenerateFromValues(millstats.FromValuesParams{
		Randomness: false,
	})

	if generator.HasErr() {
		cwd, err := os.Getwd()
		return races.Template{}, logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(err, generator.GetErr()),
		})
	}

	/*
		NOTE: Generated bonuses may be negative.
		First of all, this is intentional and follows the following logic:
		`a Troll is stupid and therefore receives an Intelligence penalty.``

		Then, this should not be handled by the process itself, but by the template
		selected with values["encounter"].

		These are configurable in the /data/systems/rules/oracles/encounter directory.
	*/
	raw := races.Template{
		Name:        values["name"],
		Description: values["description"],
		Playable:    false,
		Bonus:       generator.GetModifiers(),
	}

	return raw, nil
}

// - Helper --//
func (tb *BeastiaryTable) assertEntry(entry map[string]string) {
	tb.ClearErr()
	if entry["name"] == "" {
		tb.SetErr(logs.Error("error.required", map[string]any{
			"Subject":  "race",
			"Property": "property.name",
		}))
	}

	if entry["description"] == "" {

		tb.SetErr(logs.Error("error.required", map[string]any{
			"Subject":  "race",
			"Property": "property.description",
		}))
	}
}

func (tb *BeastiaryTable) Ensure() {
	ensureEmbeddedTable(i19n.Localize(tablebeast), &tb.Data)
	tb.Entries = superChargeEntries(tb.Entries)

	// Entries saved before IDs existed are loaded with a zero ID.
	for i := range tb.Entries {
		if tb.Entries[i].ID.IsZero() {
			tb.Entries[i].ID = id.New()
		}
	}

	tb.syncWithConfig()
}

func (tb *BeastiaryTable) syncWithConfig() error {
	var errs []error

	for _, race := range tb.Entries {
		if err := races.AddInConfig(races.Template{
			Name:        race.Race.GetName(),
			Description: race.Race.GetDescription(),
			Playable:    race.Race.IsPlayable(),
			Bonus:       race.Race.GetBonus(),
		}); err != nil {
			errs = append(errs,
				logs.CeaseError("error.saving:race", map[string]any{
					"Subject": "race",
					"Error":   err,
					"Name":    race.Race.GetName(),
				}),
			)
		}
	}

	return errors.Join(errs...)
}

func superChargeEntries(entries []beastentry) []beastentry {
	known := makeKnowEntry(entries)
	result := makeMissingEntries(known, entries)

	return result
}

func makeKnowEntry(entries []beastentry) plumtypes.Set[string] {
	set := plumtypes.NewSet[string]()

	for _, entry := range entries {
		if entry.Race == nil {
			continue
		}
		set.Add(entry.Race.GetName())
	}

	return set
}

func makeMissingEntries(known plumtypes.Set[string], entries []beastentry) []beastentry {
	all := append([]beastentry(nil), entries...)

	for _, r := range races.List() {
		if known.Has(r.GetName()) {
			continue
		}

		known.Add(r.GetName())
		all = append(all, beastentry{
			ID:        id.New(),
			Timestamp: time.Now().UTC(),
			Race:      &r,
		})
	}

	return all
}
