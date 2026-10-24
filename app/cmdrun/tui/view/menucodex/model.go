package menucodex

import (
	"solopg/app/cmdrun/domain/gameplay/codex"
	codexform "solopg/app/cmdrun/domain/gameplay/codex/form"
	codextable "solopg/app/cmdrun/domain/gameplay/codex/table"
	"solopg/app/cmdrun/tui/shortcut"
	"solopg/app/shared/components/form"
	formpage "solopg/app/shared/components/form/page"
	"solopg/app/shared/components/models"
	"solopg/app/shared/components/navigation"
	"solopg/app/shared/services/i19n"

	"solopg/app/cmdrun/types/size"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

func NewSection(focused bool) *codexMenu {
	return &codexMenu{
		FocusableSection: navigation.NewFocusableSection(navigation.TemplateFocusableSection[*codexMenuItem]{
			ID:      "codex",
			Options: makeCodexList(),
			Focused: focused,
		},
		)}
}

type CodexMenuParams struct {
	Size  size.Size
	Codex *codex.Codex
}

type codexMenu struct {
	*navigation.FocusableSection
	isSubviewOpen bool
}

// Getters & Setters

// The Codex selector remains visible in the side panel; the selected page is
// rendered separately by the gameboard state machine.
func (m *codexMenu) GetView() string {
	return m.FocusableSection.GetView()
}

func (m *codexMenu) GetFooter() []string {
	if !m.isSubviewOpen {
		return []string{shortcut.OPEN.String()}
	}

	return nil
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
	switch msg := msg.(type) {
	case form.Close:
		m.isSubviewOpen = false
	case form.Open:
		m.isSubviewOpen = true
	case tea.KeyPressMsg:
		switch plumtypes.Key(msg.String()) {
		case plumtypes.CMD_CTRL_ENTER:
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
		"codex.locations": newpage(c.Locations, codexform.GetLocationForm),
		"codex.npcs":      newpage(c.Npcs, codexform.GetNpcForm),
		"codex.beastiary": newpage(c.Beastiary, codexform.GetBeastiaryForm),
		"codex.objects":   newpage(c.Objects, codexform.GetObjectForm),
		"codex.quest":     newpage(c.Quest, codexform.GetQuestForm),
	}
}

func newpage(table codextable.Table, makeForm form.Makefunc) *formpage.Page {
	return formpage.New(formpage.Template{
		Items: func() []formpage.Item {

			entries := table.FormEntries()
			items := make([]formpage.Item, 0, len(entries))

			for _, entry := range entries {
				items = append(items, formpage.Item{
					ID:          entry.ID,
					Name:        i19n.Localize(entry.GetName()),
					Description: i19n.Localize(entry.GetDescription()),
					Values:      entry.Values,
					Locked:      entry.Locked,
				})
			}

			return items
		},
		NewForm: makeForm,
		Save:    table.SaveFromMap,
	})
}
