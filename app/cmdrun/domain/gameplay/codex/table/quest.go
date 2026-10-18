package codextable

import (
	"errors"
	"os"
	"solopg/app/cmdrun/domain/card/attributes/priority"
	"solopg/app/cmdrun/domain/card/attributes/rank"
	"solopg/app/cmdrun/domain/card/quest"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"strconv"
	"time"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// - Table --//
type QuestTable struct {
	*Data[questentry] `bson:",inline"`
}
type questentry struct {
	Timestamp time.Time    `bson:"timestamp"`
	Quest     *quest.Quest `bson:"quest"`
}

var tableobjectives = "codex.quest"

func NewQuests(entries ...questentry) *QuestTable {
	return &QuestTable{
		Data: new(i19n.Localize(tableobjectives), entries),
	}
}

// Getters & Setters

// Methods
func (tb *QuestTable) Add(quest *quest.Quest) {
	tb.Entries = append(tb.Entries, questentry{
		Timestamp: time.Now().UTC(),
		Quest:     quest,
	})
}

func (tb *QuestTable) AddFromMap(values map[string]string) error {
	tb.assertEntry(values)

	if tb.HasErr() {
		cwd, cwderr := os.Getwd()
		return logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, tb.GetErr()),
		})
	}

	// TODO convet values to quest.Quest struct un fois l'assert fait

	quest, err := quest.New(quest.Template{
		Name:        values["name"],
		Description: values["description"],
		Done:        plumtransform.StringBool(values["done"]),
		Priority:    priority.MakeDefault(values["priority"]),
		Rank:        rank.AssertWithDefault(values["rank"]),
		Reward:      values["reward"],
	})

	if err != nil {
		return logs.Error("error.invalid:new", map[string]any{
			"Subject": "objectives",
			"Error":   err,
		})
	}

	tb.Add(quest)

	logs.SilentSuccess("codex.entry:added", map[string]any{
		"Table": tableobjectives,
		"Entry": plumtransform.ParseJson(values),
	})

	return nil
}

// MEDIUM refavtor tout ca en se basant sur object ect
func (tb *QuestTable) FormEntries() []Entry {
	entries := make([]Entry, 0, len(tb.Entries))

	for _, entry := range tb.Entries {
		if entry.Quest == nil {
			continue
		}

		name, description := entry.Quest.GetName(), entry.Quest.GetDescription()

		entries = append(entries, newEntry(name, description, map[string]string{
			"name": name, "description": description,
			"done":     strconv.FormatBool(entry.Quest.Done),
			"priority": entry.Quest.Priority.GetName(),
			"rank":     entry.Quest.Rank.String(),
			"reward":   entry.Quest.Reward,
		}, false))
	}
	return entries
}

func (tb *QuestTable) SaveFromMap(index int, values map[string]string) error {
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

	// TODO Check value Done
	// TODO Check value Priority
	// TODO Check value Reward
	// TODO Check value Rank
}

func (tb *QuestTable) Ensure() {
	ensureEmbeddedTable(i19n.Localize(tableobjectives), &tb.Data)
}

//
