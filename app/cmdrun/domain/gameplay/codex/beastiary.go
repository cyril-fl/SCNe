package codex

import (
	"errors"
	"fmt"
	"os"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/characters/races"
	"solopg/app/cmdrun/services/factory/millstats"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/utils/transform"
	"strings"
	"time"
)

// - Table --//
type BeastiaryTable struct {
	*TableData[BeastiaryEntry]
}

type BeastiaryEntry struct {
	Timestamp time.Time   `bson:"timestamp"`
	Race      *races.Race `bson:"race"`
}

var tablebeast = "codex.beastiary"

func NewBeastiaryTable(entries []BeastiaryEntry) *BeastiaryTable {
	entries = superChargeEntries(entries)

	return &BeastiaryTable{
		TableData: newTable(i19n.Localize(tablebeast), entries),
	}
}

// - Methods --//
func (tb *BeastiaryTable) Add(race *races.Race) {
	tb.Entries = append(tb.Entries, BeastiaryEntry{
		Timestamp: time.Now().UTC(),
		Race:      race,
	})
}

func (tb *BeastiaryTable) AddFromMappedValues(values map[string]string) error {
	if err := tb.assertEntry(values); err != nil {
		return err
	}

	generator := millstats.NewWithValues(values)
	generator.GenerateFromValues(millstats.FromValuesParams{
		Randomness: false,
	})

	if generator.HasErr() {
		cwd, err := os.Getwd()
		return logs.Error("error.unexpected", map[string]any{
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
		"Entry": transform.ParseJson(race),
	})

	return nil
}

// HIGH check le code legacy et le code dead avec la nouvelle version
func (tb *BeastiaryTable) Summaries() []string {
	if len(tb.Entries) == 0 {
		return []string{i19n.Localize("codex.entry:empty")}
	}

	content := strings.Builder{}
	for _, entry := range tb.Entries {
		if entry.Race == nil {
			content.WriteString(i19n.Localize("codex.codex.unknown:entry"))
		} else {
			fmt.Fprintf(&content, "%s\n", i19n.Localize("codex.beastiary:entry", map[string]any{
				"Name":        entry.Race.GetName(),
				"Description": entry.Race.GetDescription(),
			}))

			if len(entry.Race.GetBonus()) > 0 {
				modifiersToStats := stats.New(entry.Race.GetBonus())
				fmt.Fprintf(&content, "%s\n", i19n.Localize("stats.entry", modifiersToStats.MappedString()))
			}
		}
		content.WriteString("\n\n")
	}

	return []string{content.String()}
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
			values[modifier.Stat.String()] = transform.ParseAsString(modifier.Value)
		}
		entries = append(entries, Entry{
			Title:       name,
			Description: description,
			Values:      values,
			Locked:      races.IsLocked(name),
		})
	}
	return entries
}

func (tb *BeastiaryTable) SaveFromMappedValues(index int, values map[string]string) error {
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

// - Helper --//
func (tb *BeastiaryTable) assertEntry(entry map[string]string) error {
	var err []error

	if entry["name"] == "" {

		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "race",
			"Property": "property.name",
		}))
	}

	if entry["description"] == "" {

		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "race",
			"Property": "property.description",
		}))
	}

	if len(err) > 0 {
		return tb.formatAssertErrors(err)
	}

	return nil
}

func (tb *BeastiaryTable) Ensure() {
	ensureEmbeddedTable(i19n.Localize(tablebeast), &tb.TableData)
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

func superChargeEntries(entries []BeastiaryEntry) []BeastiaryEntry {
	preloaded := getPreloadBeastiaryEntries()

	return append(preloaded, entries...)
}

func getPreloadBeastiaryEntries() []BeastiaryEntry {
	m := make(map[string]BeastiaryEntry)

	for _, entry := range races.List() {
		m[entry.Name] = BeastiaryEntry{
			Timestamp: time.Now().UTC(),
			Race:      &entry,
		}
	}

	slice := make([]BeastiaryEntry, 0, len(m))
	for value := range m {
		slice = append(slice, m[value])
	}

	return slice
}
