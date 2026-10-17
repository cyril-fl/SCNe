package menucodex

import (
	"solopg/app/cmdrun/domain/gameplay/codex"
	"solopg/app/shared/components/form"
	formpage "solopg/app/shared/components/form/page"
	"solopg/app/shared/components/models"
	"solopg/app/shared/components/navigation"
	"solopg/app/shared/services/i19n"
	sharedtui "solopg/app/shared/tui"

	"solopg/app/cmdrun/types/size"

	tea "charm.land/bubbletea/v2"
)

func NewSection(focused bool) *codexMenu {
	return &codexMenu{
		Section: navigation.NewSection(navigation.TemplateSection[*codexMenuItem]{
			ID: "codex", Options: makeCodexList(), Focused: focused,
		},
	)}
}

type CodexMenuParams struct {
	Size  size.Size
	Codex *codex.Codex
}

type codexMenu struct{ 
	*navigation.Section 
}


// Getters & Setters

// The Codex selector remains visible in the side panel; the selected page is
// rendered separately by the gameboard state machine.
func (m *codexMenu) GetView() string { return m.Section.GetView() }

func (m *codexMenu) GetFooter() []string {
	return []string{i19n.Localize("cmd.alt+enter:open")}
}

func (m *codexMenu) getSelectedPage() *codexMenuItem {
	if item, ok := m.GetList().SelectedItem().(models.Item[*codexMenuItem]); ok {
		return item.Value()
	}
	return nil
}

// Methods

// The codex side menu is a category selector. The selected category is opened
// in the gameboard state machine, not rendered in this side menu.
func (m *codexMenu) HandleUpdate(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch sharedtui.EventKey(key.String()) {
		case sharedtui.CMD_ALT_ENTER:
			if selected := m.getSelectedPage(); selected != nil {
				return form.OpenPage(selected.id)
			}
		}
	}
	return nil
}

// Helpers

// NewFormPages creates the pages owned by the gameboard state machine. The
// side menu only selects one of these page IDs.
func NewFormPages(c *codex.Codex) map[string]*formpage.Page {
	return map[string]*formpage.Page{
		"codex.locations": newFormPage(c.LocationsTable, codex.GetLocationForm),
		"codex.npcs":      newFormPage(c.NpcsTable, codex.GetNpcForm),
		"codex.beastiary": newFormPage(c.BeastiaryTable, codex.GetBeastiaryForm),
		"codex.objects":   newFormPage(c.ObjectsTable, codex.GetObjectForm),
		"codex.quest":     newFormPage(c.QuestTable, codex.GetQuestForm),
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

