package depracatedmenu

import (
	"solopg/app/cmdrun/domain/gameplay/hint"
	"solopg/app/shared/components/models"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type himenu struct {
	id      string
	focused bool
	list    list.Model
}


// Handlers
func (m *himenu)HiHandleUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch sharedtui.EventKey(msg.String()) {
		case sharedtui.CMD_SHIFT_ENTER, sharedtui.CMD_ALT_ENTER:
			return m.hihandleKeyShiftEnter()
		default:
			return nil
		}
	default:
		return nil
	}
}

func (m *himenu) hihandleKeyShiftEnter() tea.Cmd {
	selected, ok := m.list.SelectedItem().(models.Item[hint.Hint])
	if !ok {
		return func() tea.Msg {

			return sharedtui.SendErrorMsg(logs.Error("error.invalid", map[string]any{
				"Subject": i19n.Localize("hint"),
				"Value":   selected,
			}))
		}
	}

	spark := selected.Value()

	return func() tea.Msg {
		return HiMsg{
			Result: hint.Roll(spark),
		}
	}
}

// Messages
type HiMsg struct {
	Result []string
}
