package depracatedmenu

import (
	"solopg/app/cmdrun/domain/gameplay/dice"
	"solopg/app/shared/components/models"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// **DICE** //
type menu struct {
	id      string
	focused bool
	list    list.Model
}

func (m *menu) HandleUpdate(msg tea.Msg)  tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch sharedtui.EventKey(msg.String()) {
		case sharedtui.CMD_SHIFT_ENTER, sharedtui.CMD_ALT_ENTER:
			return  m.handleKeyShiftEnter()
		default:
			return  nil
		}
	default:
		return nil
	}
}

func (m *menu) handleKeyShiftEnter() tea.Cmd {
	selected, ok := m.list.SelectedItem().(models.Item[dice.Dice])
	if !ok {
		return func() tea.Msg {

			return sharedtui.SendErrorMsg(logs.Error("error.invalid", map[string]any{
				"Subject": i19n.Localize("dice"),
				"Value":   selected,
			}))
		}
	}

	dice := selected.Value()
	return func() tea.Msg {
		return Msg{
			Dice:  dice.GetName(),
			Value: dice.Roll(),
		}
	}
}

// Messages
type Msg struct {
	Dice  string
	Value int
}
