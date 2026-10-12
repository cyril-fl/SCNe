package depracatedmenu

import (
	"solopg/app/cmdrun/domain/gameplay/oracle"
	"solopg/app/shared/components/models"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// **ORACLES** //
type oracles struct {
	id      string
	focused bool
	list    list.Model
}

// Methods
func (m *oracles) orHandleUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch sharedtui.EventKey(msg.String()) {
		case sharedtui.CMD_SHIFT_ENTER, sharedtui.CMD_ALT_ENTER:
			return  m.orhandleKeyShiftEnter()
		default:
			return nil
			// return params.Delegate(msg)
		}
	default:
		return nil
		// return params.Delegate(msg)
	}
}

func (m *oracles) orhandleKeyShiftEnter() tea.Cmd {
	selected, ok := m.list.SelectedItem().(models.Item[oracle.Oracle])

	if !ok {
		return func() tea.Msg {

			return sharedtui.SendErrorMsg(logs.Error("error.invalid", map[string]any{
				"Subject": "oracle",
				"Value":   selected,
			}))
		}
	}

	rules := selected.Value()

	rollResult, err := oracle.Roll[any](rules)
	if err != nil {
		return func() tea.Msg {

			return sharedtui.SendErrorMsg(logs.Error("error.unexpected:action:oracle", map[string]any{
				"Action": "unexpected:action.roll",
				"Error":  err,
			}))
		}
	}

	return func() tea.Msg {
		return orMsg{
			Result: rollResult,
		}
	}
}

// Messages
type orMsg struct {
	Result *oracle.Result[any]
}

// Helpers
