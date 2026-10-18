package codextable

import (
	plumprimitive "github.com/cyril-fl/plum/primitive"
)

// - Table - //
type Table interface {
	AddFromMap(values map[string]string) error
	SaveFromMap(index int, values map[string]string) error
	FormEntries() []Entry
	Ensure()
}

// Data is embedded by each Codex table so its entries are part of the table's
// BSON document. It must remain exported because the MongoDB BSON encoder
// ignores anonymous fields whose type is unexported.
type Data[E any] struct {
	plumprimitive.Fallible

	name    string
	Entries []E `bson:"entries"`
}

func new[E any](name string, entries []E) *Data[E] {
	return &Data[E]{
		name:    name,
		Entries: entries,
	}
}

// Getters a Setters

// Methods

// Helpers
func ensureEmbeddedTable[E any](name string, table **Data[E]) {
	if *table == nil {
		*table = new(name, []E{})
	}
}

// - Entry - //

// Entry
type Entry struct {
	plumprimitive.Descriptible

	Values map[string]string
	Locked bool
}

// Getters and Setters

// Methodss
// TODO sceptiqquesur le string map
func newEntry(name, description string, values map[string]string, locked bool) Entry {
	return Entry{
		Descriptible: plumprimitive.Describe(plumprimitive.TemplateDescriptible{
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
