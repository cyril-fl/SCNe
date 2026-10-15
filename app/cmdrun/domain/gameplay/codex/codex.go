package codex

import (
	"errors"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/types/primitive"
)

// - Codex - //
type Codex struct {
	NpcsTable      *NpcsTable
	BeastiaryTable *BeastiaryTable
	LocationsTable *LocationsTable
	ObjectsTable   *ObjectsTable
	ObjectifsTable *QuestTable
}

func New() *Codex {
	return &Codex{
		NpcsTable:      NewNpcsTable([]NpcsEntry{}),
		BeastiaryTable: NewBeastiaryTable([]BeastiaryEntry{}),
		LocationsTable: NewLocationsTable([]LocationsEntry{}),
		ObjectsTable:   NewObjectsTable([]ObjectEntry{}),
		ObjectifsTable: NewQuestTable([]QuestEntry{}),
	}
}

// EnsureInitialized repairs a Codex loaded from storage where nested tables may be nil.
func (c *Codex) EnsureInitialized() {
	if c.NpcsTable == nil {
		c.NpcsTable = NewNpcsTable(nil)
	}
	if c.BeastiaryTable == nil {
		c.BeastiaryTable = NewBeastiaryTable(nil)
	}
	if c.LocationsTable == nil {
		c.LocationsTable = NewLocationsTable(nil)
	}
	if c.ObjectsTable == nil {
		c.ObjectsTable = NewObjectsTable(nil)
	}
	if c.ObjectifsTable == nil {
		c.ObjectifsTable = NewQuestTable(nil)
	}

	c.NpcsTable.Ensure()
	c.BeastiaryTable.Ensure()
	c.LocationsTable.Ensure()
	c.ObjectsTable.Ensure()
	c.ObjectifsTable.Ensure()
}

// - Codex table - //
type Table interface {
	AddFromMappedValues(values map[string]string) error
	SaveFromMappedValues(index int, values map[string]string) error
	FormEntries() []Entry
	Ensure()
}

type Entry struct {
	primitive.Descriptible

	Values map[string]string
	Locked bool
}

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

func ensureEmbeddedTable[E any](name string, table **TableData[E]) {
	if *table == nil {
		*table = newTable(name, []E{})
	}
}

func (t TableData[any]) formatAssertErrors(err []error) error {
	return errors.Join(
		logs.Error("error.invalid", map[string]interface{}{
			"Subject": t.name,
			"Value":   err,
		}),

		errors.Join(err...),
	)
}
