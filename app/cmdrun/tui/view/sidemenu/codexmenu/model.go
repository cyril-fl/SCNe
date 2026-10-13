package codexmenu

import (
	"solopg/app/cmdrun/domain/gameplay/codex"
	"solopg/app/shared/components/models"
	"solopg/app/shared/components/navigation"
	sharedtui "solopg/app/shared/tui"

	"solopg/app/cmdrun/types/direction"
	"solopg/app/cmdrun/types/size"

	tea "charm.land/bubbletea/v2"
)

func NewSideMenu(params CodexMenuParams, focused bool) *codexMenu {
	return &codexMenu{Menu: navigation.NewMenu(navigation.TemplateMenu[*codexMenuItem]{
		ID: "codex", Options: makeCodexList(params), Focused: focused,
	})}
}

type CodexMenuParams struct {
	Size  size.Size
	Codex *codex.Codex
}

type codexMenu struct{ *navigation.Menu }

// The codex side menu is a category selector. The selected category is opened
// in the gameboard state machine, not rendered in this side menu.
func (m *codexMenu) IsOpen() bool                             { return false }
func (m *codexMenu) SetOpen(bool)                             {}
func (m *codexMenu) HandleDirectionInput(direction.Direction) {}

func (m *codexMenu) HandleUpdate(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch sharedtui.EventKey(key.String()) {
		case sharedtui.CMD_ALT_ENTER:
			if selected := m.getSelectedPage(); selected != nil {
				return sendOpenPage(selected.id)
			}
		}
	}
	return nil
}

type codexMenuItem struct{ id string }

func (m *codexMenuItem) GetName() string { return m.id }

func (m *codexMenu) getSelectedPage() *codexMenuItem {
	if item, ok := m.GetList().SelectedItem().(models.Item[*codexMenuItem]); ok {
		return item.Value()
	}
	return nil
}

type OpenPageMsg struct{ ID string }

func sendOpenPage(id string) tea.Cmd {
	return func() tea.Msg { return OpenPageMsg{ID: id} }
}
