package codextable

import (
	"errors"
	"os"
	"solopg/app/cmdrun/domain/card/quest"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/utils/transform"
	"time"
)

// - Table --//
type QuestTable struct {
	*TableData[QuestEntry]
}
type QuestEntry struct {
	Timestamp time.Time    `bson:"timestamp"`
	Quest     *quest.Quest `bson:"quest"`
}

var tableobjectives = "codex.quest"

func NewQuestTable(entries []QuestEntry) Table {
	return &QuestTable{
		TableData: newTable(i19n.Localize(tableobjectives), entries),
	}
}

// Getters & Setters

// Methods
func (o *QuestTable) Add(quest *quest.Quest) {
	o.Entries = append(o.Entries, QuestEntry{
		Timestamp: time.Now().UTC(),
		Quest:     quest,
	})
}

func (o *QuestTable) AddFromMappedValues(values map[string]string) error {
	if err := o.assertEntry(values); err != nil {
		cwd, cwderr := os.Getwd()
		return logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, err),
		})
	}

	quest, err := quest.New(quest.Template{
		Name:        values["name"],
		Description: values["description"],
	})

	if err != nil {
		return logs.Error("error.invalid:new", map[string]any{
			"Subject": "objectives",
			"Error":   err,
		})
	}

	o.Add(quest)

	logs.SilentSuccess("codex.entry:added", map[string]any{
		"Table": tableobjectives,
		"Entry": transform.ParseJson(values),
	})

	return nil
}

// MEDIUM refavtor tout ca en se basant sur object ect
func (o *QuestTable) FormEntries() []Entry {
	entries := make([]Entry, 0, len(o.Entries))

	for _, entry := range o.Entries {
		if entry.Quest == nil {
			continue
		}

		name, description := entry.Quest.GetName(), entry.Quest.GetDescription()

		entries = append(entries, newEntry(name, description, map[string]string{
			"name": name, "description": description,
		}, false))
	}
	return entries
}

func (o *QuestTable) SaveFromMappedValues(index int, values map[string]string) error {
	if index < 0 {
		return o.AddFromMappedValues(values)
	}
	oldLen := len(o.Entries)
	if err := o.AddFromMappedValues(values); err != nil {
		return err
	}
	if index < oldLen {
		o.Entries[index] = o.Entries[oldLen]
		o.Entries = o.Entries[:oldLen]
	}
	return nil
}

func (o *QuestTable) FindEntryByName(name string) *Entry {
	return findEntryByName(o.FormEntries(), name)
}

func (o *QuestTable) assertEntry(values map[string]string) error {
	var err []error

	if values["name"] == "" {
		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "objectives",
			"Property": "property.name",
		}))
	}
	if values["description"] == "" {

		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "objectives",
			"Property": "property.description",
		}))
	}

	if len(err) > 0 {
		return o.formatAssertErrors(err)
	}

	return nil
}

func (o *QuestTable) Ensure() {
	ensureEmbeddedTable(i19n.Localize(tableobjectives), &o.TableData)
}

//
