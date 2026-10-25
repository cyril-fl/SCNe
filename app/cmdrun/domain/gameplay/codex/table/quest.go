package codextable

import (
	"errors"
	"os"
	"scne/app/cmdrun/domain/card/attributes/priority"
	"scne/app/cmdrun/domain/card/attributes/rank"
	"scne/app/cmdrun/domain/card/quest"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/logs"
	"slices"
	"strconv"
	"time"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// - Table --//

// QuestTable lists the quests of the campaign.
type QuestTable struct {
	*Data[questentry] `bson:",inline"`
}
type questentry struct {
	ID        id.ID        `bson:"id"`
	Timestamp time.Time    `bson:"timestamp"`
	Quest     *quest.Quest `bson:"quest"`
}

var tableobjectives = "codex.quest"

// NewQuests returns a quest table with entries.
func NewQuests(entries ...questentry) *QuestTable {
	return &QuestTable{
		Data: new(i19n.Localize(tableobjectives), entries),
	}
}

// Getters & Setters

// Methods

// Add appends a new quest entry.
func (tb *QuestTable) Add(quest *quest.Quest) {
	tb.Entries = append(tb.Entries, questentry{
		ID:        id.New(),
		Timestamp: time.Now().UTC(),
		Quest:     quest,
	})
}

// AddFromMap creates a quest from form values and appends it.
func (tb *QuestTable) AddFromMap(values map[string]string) error {
	quest, err := tb.questFromMap(values)
	if err != nil {
		return err
	}

	tb.Add(quest)

	logs.SilentSuccess("codex.entry:added", map[string]any{
		"Table": tableobjectives,
		"Entry": plumtransform.ParseJson(values),
	})

	return nil
}

func (tb *QuestTable) questFromMap(values map[string]string) (*quest.Quest, error) {
	tb.assertEntry(values)

	if tb.HasErr() {
		cwd, cwderr := os.Getwd()
		return nil, logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, tb.GetErr()),
		})
	}

	quest, err := quest.New(quest.Template{
		Name:        values["name"],
		Description: values["description"],
		Done:        plumtransform.StringBool(values["done"]),
		Priority:    priority.MakeDefault(values["priority"]),
		Rank:        rank.AssertWithDefault(values["rank"]),
		Reward:      values["reward"],
	})

	if err != nil {
		return nil, logs.Error("error.invalid:new", map[string]any{
			"Subject": "objectives",
			"Error":   err,
		})
	}

	return quest, nil
}

// FormEntries returns the entries as form page rows.
func (tb *QuestTable) FormEntries() []Entry {
	entries := make([]Entry, 0, len(tb.Entries))

	for _, entry := range tb.Entries {
		if entry.Quest == nil {
			continue
		}

		name, description := entry.Quest.GetName(), entry.Quest.GetDescription()

		entries = append(entries, newEntry(entry.ID, name, description, map[string]string{
			"name":        name,
			"description": description,
			"done":        strconv.FormatBool(entry.Quest.Done),
			"priority":    entry.Quest.Priority.GetName(),
			"rank":        entry.Quest.Rank.String(),
			"reward":      entry.Quest.Reward,
		}, false))
	}
	return entries
}

// SaveFromMap edits the entry with entryID, or adds one when it is zero or unknown.
func (tb *QuestTable) SaveFromMap(entryID id.ID, values map[string]string) error {
	index := slices.IndexFunc(tb.Entries, func(entry questentry) bool { return entry.ID == entryID })
	if entryID.IsZero() || index < 0 {
		return tb.AddFromMap(values)
	}
	return tb.editFromMap(index, values)
}

func (tb *QuestTable) editFromMap(index int, values map[string]string) error {
	edited, err := tb.questFromMap(values)
	if err != nil {
		return err
	}

	tb.Entries[index].Quest = edited
	return nil
}

func (tb *QuestTable) assertEntry(values map[string]string) {
	tb.ClearErr()

	if values["name"] == "" {
		tb.SetErr(logs.Error("error.required", map[string]any{
			"Subject":  "objectives",
			"Property": "property.name",
		}))
	}
	if values["description"] == "" {
		tb.SetErr(logs.Error("error.required", map[string]any{
			"Subject":  "objectives",
			"Property": "property.description",
		}))
	}
}

// Ensure initializes the table when missing.
func (tb *QuestTable) Ensure() {
	ensureEmbeddedTable(i19n.Localize(tableobjectives), &tb.Data)
}

//
