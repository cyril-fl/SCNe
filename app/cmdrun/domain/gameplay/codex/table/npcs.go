package codextable

import (
	"errors"
	"os"
	"slices"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/characters"
	"solopg/app/cmdrun/domain/card/characters/classes"
	"solopg/app/cmdrun/domain/card/characters/races"
	"solopg/app/cmdrun/services/factory/millstats"
	"solopg/app/cmdrun/types/id"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"strconv"

	"time"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// - Table --//
type NpcsTable struct {
	*Data[npcsentry] `bson:",inline"`
}
type npcsentry struct {
	ID        id.ID                 `bson:"id"`
	Timestamp time.Time             `bson:"timestamp"`
	Character *characters.Character `bson:"character"`
}

var tablenpcs = "codex.npcs"

func NewNpcs(entries ...npcsentry) *NpcsTable {
	return &NpcsTable{
		Data: new(i19n.Localize(tablenpcs), entries),
	}
}

// Getters & Setters

// Methods
func (tb *NpcsTable) Add(character *characters.Character) {
	tb.Entries = append(tb.Entries, npcsentry{
		ID:        id.New(),
		Timestamp: time.Now().UTC(),
		Character: character,
	})
}

func (tb *NpcsTable) AddFromMap(values map[string]string) error {
	if err := tb.checkEntry(values); err != nil {
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
	character, err := newCharacterFromMap(values, generator.GetStats())
	if err != nil {
		return err
	}

	tb.Add(character)

	logs.SilentSuccess("codex.entry:added", map[string]any{
		"Table": tablenpcs,
		"Entry": plumtransform.ParseJson(character),
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
			values[stat] = plumtransform.ParseAsString(value)
		}
		entries = append(entries, newEntry(entry.ID,
			character.GetName(),
			character.Race,
			values,
			characters.IsLocked(character.GetName()),
		))
	}
	return entries
}

func (tb *NpcsTable) SaveFromMap(entryID id.ID, values map[string]string) error {
	index := slices.IndexFunc(tb.Entries, func(entry npcsentry) bool { return entry.ID == entryID })
	if entryID.IsZero() || index < 0 {
		return tb.AddFromMap(values)
	}
	return tb.editFromMap(index, values)
}

// Edited values already hold final stats: only the class and race bonus
// difference is applied, never the roll.
func (tb *NpcsTable) editFromMap(index int, values map[string]string) error {
	if err := tb.checkEntry(values); err != nil {
		return err
	}

	previous := tb.Entries[index].Character
	previousBonus := getBonus(previous.Class, previous.Race)
	editedBonus := getBonus(values["class"], values["race"])

	editedStats := statsFromValues(values)
	editedStats.ApplyModifiers(negateModifiers(previousBonus))
	editedStats.ApplyModifiers(editedBonus)

	character, err := newCharacterFromMap(values, editedStats)
	if err != nil {
		return err
	}

	tb.Entries[index].Character = character
	return nil
}

func (tb *NpcsTable) checkEntry(values map[string]string) error {
	tb.assertEntry(values)

	if tb.HasErr() {
		cwd, cwderr := os.Getwd()
		return logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, tb.GetErr()),
		})
	}
	return nil
}

func (tb *NpcsTable) assertEntry(entry map[string]string) {
	tb.ClearErr()

	if entry["name"] == "" {
		tb.SetErr(logs.Error("error.required", map[string]any{
			"Subject":  "npcs",
			"Property": "property.name",
		}))
	}
	if entry["description"] == "" {
		tb.SetErr(logs.Error("error.required", map[string]any{
			"Subject":  "npcs",
			"Property": "property.description",
		}))
	}
	if !races.Assert(entry["race"]) {
		tb.SetErr(logs.Error("error.required", map[string]any{
			"Subject":  "npcs",
			"Property": "property.race",
		}))
	}
	if !classes.Assert(entry["class"]) {
		tb.SetErr(logs.Error("error.required", map[string]any{
			"Subject":  "npcs",
			"Property": "property.class",
		}))
	}
}

func (tb *NpcsTable) Ensure() {
	ensureEmbeddedTable(i19n.Localize(tablenpcs), &tb.Data)

	// Entries saved before IDs existed are loaded with a zero ID.
	for i := range tb.Entries {
		if tb.Entries[i].ID.IsZero() {
			tb.Entries[i].ID = id.New()
		}
	}
}

// Helpers
func newCharacterFromMap(values map[string]string, characterStats stats.Stats) (*characters.Character, error) {
	character, err := characters.New(characters.Template{
		Name:        values["name"],
		Description: values["description"],
		Rarity:      rarity.Default(),
		Race:        values["race"],
		Class:       values["class"],
		Stats:       characterStats,
	})

	if err != nil {
		return nil, logs.Error("error.invalid:new", map[string]any{
			"Subject": "npcs",
			"Error":   err,
		})
	}
	return character, nil
}

func statsFromValues(values map[string]string) stats.Stats {
	result := make(stats.Stats)
	for _, stat := range stats.List() {
		value, _ := strconv.Atoi(values[stat.String()])
		result[stat] = value
	}
	return result
}

func getBonus(class, race string) []stats.Modifier {
	var bonus []stats.Modifier
	if c := classes.FindByName(class); c != nil {
		bonus = append(bonus, c.GetBonus()...)
	}
	if r := races.FindByName(race); r != nil {
		bonus = append(bonus, r.GetBonus()...)
	}
	return bonus
}

func negateModifiers(modifiers []stats.Modifier) []stats.Modifier {
	negated := make([]stats.Modifier, 0, len(modifiers))
	for _, modifier := range modifiers {
		modifier.Value = -modifier.Value
		negated = append(negated, modifier)
	}
	return negated
}
