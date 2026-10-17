package codextable

import (
	"errors"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/types/primitive"
)

// - Table - //
type Table interface {
	AddFromMappedValues(values map[string]string) error
	SaveFromMappedValues(index int, values map[string]string) error
	FormEntries() []Entry
	FindEntryByName(name string) *Entry
	Ensure()
}

type TableData[E any] struct {
	name    string
	Entries []E `bson:"entries"`
}

func newTable[E any](name string, entries []E) *TableData[E] {
	return &TableData[E]{
		name:    name,
		Entries: entries,
	}
}

// Getters a Setters

// Methods
func (t TableData[any]) formatAssertErrors(err []error) error {
	return errors.Join(
		logs.Error("error.invalid", map[string]interface{}{
			"Subject": t.name,
			"Value":   err,
		}),

		errors.Join(err...),
	)
}

// Helpers
func ensureEmbeddedTable[E any](name string, table **TableData[E]) {
	if *table == nil {
		*table = newTable(name, []E{})
	}
}

// - Entry - //

// Entry
type Entry struct {
	primitive.Descriptible

	Values map[string]string
	Locked bool
}

// Getters and Setters

// Methodss
// TODO sceptiqquesur le string map
func newEntry(name, description string, values map[string]string, locked bool) Entry {
	return Entry{
		Descriptible: primitive.Describe(primitive.TemplateDescriptible{
			Name:        name,
			Description: description,
		}),
		Values: values,
		Locked: locked,
	}
}

// Helpers
func findEntryByName(entries []Entry, name string) *Entry {
	for i := range entries {
		if entries[i].GetName() == name {
			return &entries[i]
		}
	}
	return nil
}
