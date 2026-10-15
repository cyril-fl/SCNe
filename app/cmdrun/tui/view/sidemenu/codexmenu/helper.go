package codexmenu

import (
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/characters/classes"
	"solopg/app/cmdrun/domain/card/characters/races"
	"solopg/app/cmdrun/domain/gameplay/codex"
	"solopg/app/cmdrun/domain/gameplay/oracle"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/form/field"
	"solopg/app/shared/components/formpage"
	"solopg/app/shared/components/models"
	"solopg/app/shared/services/i19n"
	"solopg/config"
)

func makeCodexList(params CodexMenuParams) []*codexMenuItem {
	return []*codexMenuItem{
		{id: "codex.locations"},
		{id: "codex.npcs"},
		{id: "codex.beastiary"},
		{id: "codex.objects"},
		{id: "codex.objectives"},
	}
}

// NewFormPages creates the pages owned by the gameboard state machine. The
// side menu only selects one of these page IDs.
func NewFormPages(c *codex.Codex) map[string]*formpage.FormPage {
	return map[string]*formpage.FormPage{
		"codex.locations":  newFormPage(c.LocationsTable, getLocationForm),
		"codex.npcs":       newFormPage(c.NpcsTable, getNpcForm),
		"codex.beastiary":  newFormPage(c.BeastiaryTable, getBeastiaryForm),
		"codex.objects":    newFormPage(c.ObjectsTable, getObjectForm),
		"codex.objectives": newFormPage(c.ObjectifsTable, getObjectifForm),
	}
}

func newFormPage(table codex.Table, makeForm func(map[string]string) *form.Form) *formpage.FormPage {
	return formpage.New(formpage.Template{
		Items: func() []formpage.Item {
			entries := table.FormEntries()
			items := make([]formpage.Item, 0, len(entries))
			for _, entry := range entries {
				items = append(items, formpage.Item{
					Title:       i19n.Localize(entry.GetName()),
					Description: i19n.Localize(entry.GetDescription()),
					Values:      entry.Values,
					Locked:      entry.Locked,
				})
			}
			return items
		},
		NewForm: makeForm,
		Save:    table.SaveFromMappedValues,
	})
}

func getFormRaceOptions() []models.Item[string] {
	items := make([]models.Item[string], 0)
	for _, value := range races.List() {
		items = append(items, models.NewItem(i19n.Localize(value.GetName()), "", value.GetName()))
	}
	return items
}

func getFormClassOptions() []models.Item[string] {
	items := make([]models.Item[string], 0)
	for _, value := range classes.List() {
		items = append(items, models.NewItem(i19n.Localize(value.GetName()), "", value.GetName()))
	}
	return append(items, models.NewItem("None", "", ""))
}

func getFormOracle() []models.Item[string] {
	values, err := oracle.ListFromFolder(config.Current.Documents.Folders.Encounters)
	if err != nil {
		return nil
	}
	items := make([]models.Item[string], 0, len(values))
	for _, value := range values {
		items = append(items, models.NewItem(i19n.Localize("oracles."+value), "", value))
	}
	return items
}

func getFormStatsFields(values map[string]string) []form.Field {
	fields := make([]form.Field, 0)
	for _, value := range stats.List() {
		fields = append(fields, field.TextField(field.TextTemplate[string]{
			ID:           value.String(),
			Label:        i19n.Localize(value.String()),
			Defaultvalue: values[value.String()],
		}))
	}
	return fields
}
