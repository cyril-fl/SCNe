package codextable

import (
	"errors"
	"os"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/characters"
	"solopg/app/cmdrun/domain/card/characters/classes"
	"solopg/app/cmdrun/domain/card/characters/races"
	"solopg/app/cmdrun/services/factory/millstats"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/utils/transform"

	"time"
)

// - Table --//
type NpcsTable struct {
	*TableData[NpcsEntry]
}
type NpcsEntry struct {
	Timestamp time.Time             `bson:"timestamp"`
	Character *characters.Character `bson:"character"`
}

var tablenpcs = "codex.npcs"

func NewNpcsTable(entries []NpcsEntry) Table {
	return &NpcsTable{
		TableData: newTable(i19n.Localize(tablenpcs), entries),
	}
}

// Getters & Setters

// Methods
func (tb *NpcsTable) Add(character *characters.Character) {
	tb.Entries = append(tb.Entries, NpcsEntry{
		Timestamp: time.Now().UTC(),
		Character: character,
	})
}

func (tb *NpcsTable) AddFromMappedValues(values map[string]string) error {
	if err := tb.assertEntry(values); err != nil {
		return err
	}

	generator := millstats.NewWithValues(values)
	generator.GenerateFromValues(millstats.FromValuesParams{
		Randomness: true,
	})

	if generator.HasErr() {
		cwd, err := os.Getwd()
		return logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(err, generator.GetErr()),
		})
	}

	/*
		NOTE Description should be a text field that describes the character in a RP way, and includes details
		about the character's equipment, abilities, and other relevant information. Power level should
		be calculated / represented by a dice roll in an oracle, for example:
		a powerful character stats would be chosen with the legendary.
	*/
	character, err := characters.New(characters.Template{
		Name:        values["name"],
		Description: values["description"],
		Rarity:      rarity.Default(),
		Race:        values["race"],
		Class:       values["class"],
		Stats:       generator.GetStats(),
	})

	if err != nil {
		return logs.Error("error.invalid:new", map[string]any{
			"Subject": "npcs",
			"Error":   err,
		})
	}

	tb.Add(character)

	logs.SilentSuccess("codex.entry:added", map[string]any{
		"Table": tablenpcs,
		"Entry": transform.ParseJson(character),
	})

	return nil
}

func (tb *NpcsTable) FormEntries() []Entry {
	entries := make([]Entry, 0, len(tb.Entries))
	for _, entry := range tb.Entries {
		if entry.Character == nil {
			continue
		}
		character := entry.Character
		values := map[string]string{
			"name": character.GetName(), "description": character.GetDescription(),
			"class": character.Class, "race": character.Race,
		}
		for stat, value := range character.Stats.MappedString() {
			values[stat] = transform.ParseAsString(value)
		}
		entries = append(entries, newEntry(
			character.GetName(), character.Race, values,
			characters.IsLocked(character.GetName()),
		))
	}
	return entries
}

func (tb *NpcsTable) SaveFromMappedValues(index int, values map[string]string) error {
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

func (tb *NpcsTable) FindEntryByName(name string) *Entry {
	return findEntryByName(tb.FormEntries(), name)
}

func (tb *NpcsTable) assertEntry(entry map[string]string) error {
	var err []error
	if entry["name"] == "" {

		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "npcs",
			"Property": "property.name",
		}))
	}
	if entry["description"] == "" {

		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "npcs",
			"Property": "property.description",
		}))
	}
	if !races.Assert(entry["race"]) {

		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "npcs",
			"Property": "property.race",
		}))
	}
	if !classes.Assert(entry["class"]) {

		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "npcs",
			"Property": "property.class",
		}))
	}

	if len(err) > 0 {
		return tb.formatAssertErrors(err)
	}

	return nil
}

func (tb *NpcsTable) Ensure() {
	ensureEmbeddedTable(i19n.Localize(tablenpcs), &tb.TableData)
}

// Helpers
