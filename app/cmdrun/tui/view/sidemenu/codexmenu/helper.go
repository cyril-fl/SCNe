package codexmenu

import (
	"solopg/app/cmdrun/domain/gameplay/codex"
	"solopg/app/cmdrun/tui/view/formview"
	"solopg/app/shared/components/form"
	formpage "solopg/app/shared/components/form/page"

	"solopg/app/shared/services/i19n"
)

func makeCodexList() []*codexMenuItem {
	return []*codexMenuItem{
		{id: "codex.locations"},
		{id: "codex.npcs"},
		{id: "codex.beastiary"},
		{id: "codex.objects"},
		{id: "codex.quest"},
	}
}

// NewFormPages creates the pages owned by the gameboard state machine. The
// side menu only selects one of these page IDs.
func NewFormPages(c *codex.Codex) map[string]*formpage.Page {
	return map[string]*formpage.Page{
		"codex.locations": newFormPage(c.LocationsTable, formview.GetLocationForm),
		"codex.npcs":      newFormPage(c.NpcsTable, formview.GetNpcForm),
		"codex.beastiary": newFormPage(c.BeastiaryTable, formview.GetBeastiaryForm),
		"codex.objects":   newFormPage(c.ObjectsTable, formview.GetObjectForm),
		"codex.quest":     newFormPage(c.QuestTable, formview.GetQuestForm),
	}
}

func newFormPage(table codex.Table, makeForm func(map[string]string) *form.Form) *formpage.Page {
	return formpage.New(formpage.Template{
		Items: func() []formpage.Item {
			entries := table.FormEntries()
			items := make([]formpage.Item, 0, len(entries))
			for _, entry := range entries {
				items = append(items, formpage.Item{
					Name:        i19n.Localize(entry.GetName()),
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
