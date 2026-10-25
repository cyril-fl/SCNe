package menucodex

import (
	"solopg/app/cmdrun/domain/gameplay/codex"
	codexform "solopg/app/cmdrun/domain/gameplay/codex/form"
	codextable "solopg/app/cmdrun/domain/gameplay/codex/table"
	"solopg/app/cmdrun/domain/gameplay/portal"
	"solopg/app/cmdrun/domain/gameplay/session"
	"solopg/app/cmdrun/tui/shortcut"
	"solopg/app/cmdrun/types/id"
	"solopg/app/shared/components/form"
	formpage "solopg/app/shared/components/form/page"
	"solopg/app/shared/components/models"
	"solopg/app/shared/components/navigation"
	"solopg/app/shared/messages"
	"solopg/app/shared/services/i19n"

	"solopg/app/cmdrun/types/size"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

// COMBAT_PAGE is the gameboard page opened by the fight action.
const COMBAT_PAGE = "combat"

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
func NewFormPages(gameSession *session.GameSession) map[string]*formpage.Page {
	c := gameSession.State.Codex
	return map[string]*formpage.Page{
		"codex.locations": newpage(c.Locations, codexform.GetLocationForm, newTravelAction(gameSession)),
		"codex.npcs":      newpage(c.Npcs, codexform.GetNpcForm, newFightAction(gameSession)),
		"codex.beastiary": newpage(c.Beastiary, codexform.GetBeastiaryForm, nil),
		"codex.objects":   newpage(c.Objects, codexform.GetObjectForm, nil),
		"codex.quest":     newpage(c.Quest, codexform.GetQuestForm, nil),
	}
}

// newTravelAction sets the selected codex location as the session's current location.
func newTravelAction(gameSession *session.GameSession) *formpage.Action {
	return &formpage.Action{
		Label: "travel",
		Run: func(entryID id.ID) tea.Cmd {
			portal.Travel(gameSession, entryID)
			return nil
		},
	}
}

// newFightAction starts a combat against the selected codex NPC and opens the combat page.
func newFightAction(gameSession *session.GameSession) *formpage.Action {
	return &formpage.Action{
		Label: "fight",
		Run: func(entryID id.ID) tea.Cmd {
			opponent := gameSession.State.Codex.Npcs.Find(entryID)
			if err := gameSession.StartCombat(opponent); err != nil {
				return messages.Fail(err)
			}
			return form.OpenPage(COMBAT_PAGE)
		},
	}
}

func newpage(table codextable.Table, makeForm form.Makefunc, action *formpage.Action) *formpage.Page {
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
		Action:  action,
	})
}
