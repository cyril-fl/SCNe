package codextable

import (
	"solopg/app/shared/types/primitive"
)

// - Table - //
type Table interface {
	AddFromMap(values map[string]string) error
	SaveFromMap(index int, values map[string]string) error
	FormEntries() []Entry
	Ensure()
}

type data[E any] struct {
	primitive.Fallible

	name    string
	Entries []E `bson:"entries"`
}

func new[E any](name string, entries []E) *data[E] {
	return &data[E]{
		name:    name,
		Entries: entries,
	}
}

// Getters a Setters

// Methods

// Helpers
func ensureEmbeddedTable[E any](name string, table **data[E]) {
	if *table == nil {
		*table = new(name, []E{})
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
func FindEntryByName(entries []Entry, name string) *Entry {
	for i := range entries {
		if entries[i].GetName() == name {
			return &entries[i]
		}
	}
	return nil
}
