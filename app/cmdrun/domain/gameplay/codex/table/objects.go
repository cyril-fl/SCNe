package codextable

import (
	"solopg/app/cmdrun/domain/card/attributes/objectcategory"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/cmdrun/domain/card/objects"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/utils/transform"
	"time"
)

// - Table --//
/*
LOW Modifier Object table pour quelle prenne des card weapon / article en fonction de et creer la bonne factory pour aller avec
*/
type ObjectsTable struct {
	*TableData[ObjectEntry]
}

type ObjectEntry struct {
	Timestamp time.Time       `bson:"timestamp"`
	Object    *objects.Object `bson:"object"`
}

var tableobject = "codex.objects"

func NewObjectsTable(entries []ObjectEntry) Table {
	return &ObjectsTable{
		TableData: newTable(i19n.Localize(tableobject), entries),
	}
}

// Getters & Setters

// Methods
func (tb *ObjectsTable) Add(object *objects.Object) {
	tb.Entries = append(tb.Entries, ObjectEntry{
		Timestamp: time.Now().UTC(),
		Object:    object,
	})
}

func (tb *ObjectsTable) AddFromMappedValues(values map[string]string) error {
	if err := tb.assertEntry(values); err != nil {
		return err
	}

	object, err := objects.New(objects.Template{
		Name:        values["name"],
		Description: values["description"],
		Rarity:      rarity.Default(),
		Variety:     variety.AssertWithDefault(values["variety"]),
		Category:    objectcategory.AssertWithDefault(values["category"]),
	})

	if err != nil {
		return logs.Error("error.invalid:new", map[string]any{
			"Subject": "item",
			"Error":   err,
		})
	}

	tb.Add(object)

	logs.SilentSuccess("codex.entry:added", map[string]any{
		"Table": tableobject,
		"Entry": transform.ParseJson(object),
	})

	return nil
}

func (tb *ObjectsTable) FormEntries() []Entry {
	entries := make([]Entry, 0, len(tb.Entries))

	for _, entry := range tb.Entries {
		if entry.Object == nil {
			continue
		}

		name, description := entry.Object.GetName(), entry.Object.GetDescription()

		entries = append(entries, newEntry(name, description, map[string]string{
			"name": name, "description": description, "category": entry.Object.Category.String(),
		}, false))
	}
	return entries
}

func (tb *ObjectsTable) SaveFromMappedValues(index int, values map[string]string) error {
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

func (tb *ObjectsTable) FindEntryByName(name string) *Entry {
	return findEntryByName(tb.FormEntries(), name)
}

func (tb *ObjectsTable) assertEntry(entry map[string]string) error {
	var err []error

	if entry["name"] == "" {
		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "item",
			"Property": "property.name",
		}))
	}
	if entry["description"] == "" {

		err = append(err, logs.Error("error.required", map[string]any{
			"Subject":  "item",
			"Property": "property.description",
		}))
	}
	/*
		HIGH Quand tout sera fix ajouter ces champs dans le form
		if entry["rarity"] == "" {
			err = append(err, fmt.Errorf("rarity is required"))
		}
		if entry["variety"] == "" {
			err = append(err, fmt.Errorf("variety is required"))
		}
		if entry["category"] == "" {
			err = append(err, fmt.Errorf("category is required"))
		}
	*/

	if len(err) > 0 {
		return tb.formatAssertErrors(err)
	}

	return nil
}

func (tb *ObjectsTable) Ensure() {
	ensureEmbeddedTable(i19n.Localize(tableobject), &tb.TableData)
}

// Helpers
