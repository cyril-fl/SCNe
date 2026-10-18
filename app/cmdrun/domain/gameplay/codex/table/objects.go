package codextable

import (
	"errors"
	"os"
	"solopg/app/cmdrun/domain/card/attributes/objectcategory"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/cmdrun/domain/card/objects"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"time"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// - Table --//
/*
LOW Modifier Object table pour quelle prenne des card weapon / article en fonction de et creer la bonne factory pour aller avec
*/
type ObjectsTable struct {
	*data[objectentry]
}

type objectentry struct {
	Timestamp time.Time       `bson:"timestamp"`
	Object    *objects.Object `bson:"object"`
}

var tableobject = "codex.objects"

func NewObjects(entries ...objectentry) *ObjectsTable {
	return &ObjectsTable{
		data: new(i19n.Localize(tableobject), entries),
	}
}

// Getters & Setters

// Methods
func (tb *ObjectsTable) AddFromMap(values map[string]string) error {
	tb.assertEntry(values)

	if tb.HasErr() {
		cwd, cwderr := os.Getwd()
		return logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, tb.GetErr()),
		})
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
		"Entry": plumtransform.ParseJson(object),
	})

	return nil
}

func (tb *ObjectsTable) Add(object *objects.Object) {
	tb.Entries = append(tb.Entries, objectentry{
		Timestamp: time.Now().UTC(),
		Object:    object,
	})
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

func (tb *ObjectsTable) SaveFromMap(index int, values map[string]string) error {
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

func (tb *ObjectsTable) assertEntry(entry map[string]string) {
	tb.ClearErr()

	if entry["name"] == "" {
		tb.SetErr(logs.Error("error.required", map[string]any{
			"Subject":  "item",
			"Property": "property.name",
		}))
	}
	if entry["description"] == "" {

		tb.SetErr(logs.Error("error.required", map[string]any{
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
}

func (tb *ObjectsTable) Ensure() {
	ensureEmbeddedTable(i19n.Localize(tableobject), &tb.data)
}

// Helpers
