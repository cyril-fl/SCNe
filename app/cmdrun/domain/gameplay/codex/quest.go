package codex

import (
	"errors"
	"os"
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
	Timestamp   time.Time `bson:"timestamp"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
}

var tableobjectives = "codex.objectives"

func NewQuestTable(entries []QuestEntry) *QuestTable {
	return &QuestTable{
		TableData: newTable(i19n.Localize(tableobjectives), entries),
	}
}

// - Methods --//
func (o *QuestTable) Add(title, description string) {
	o.Entries = append(o.Entries, QuestEntry{
		Timestamp:   time.Now().UTC(),
		Title:       title,
		Description: description,
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

	o.Add(values["title"], values["description"])

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
		entries = append(entries, newEntry(entry.Title, entry.Description, map[string]string{
			"title": entry.Title, "description": entry.Description,
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

// - Helper --//
func (o *QuestTable) assertEntry(values map[string]string) error {
	var err []error

	if values["title"] == "" {

		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "objectives",
			"Property": "property.title",
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
