package gameboard

import (
	"solopg/app/shared/messages"
	"solopg/app/shared/services/logs"
	"time"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

func (m *model) handleSaveInput(msg messages.Saved) (*model, tea.Cmd) {
	m.SetErr(msg.Err)

	if m.HasErr() {
		err := logs.Error("error.unexpected:save", map[string]any{
			"Error": msg.Err,
		})

		m.state.Amend(err.Error())
	} else {
		logs.SilentSuccess("campaign:success", map[string]any{
			"Time": time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return m, nil
}

// - Handlers - //
// Input

func (m *model) handleKeyPress(msg tea.KeyPressMsg) (*model, tea.Cmd) {
	switch plumtypes.Key(msg.String()) {
	case plumtypes.KEY_UP, plumtypes.KEY_DOWN:
		// m.handleDirectionInput(msg)
		// default:
		// 	handleDefaultInput(m, msg)
	}
	return m.handleCommand(msg)
}

func (m *model) handleCommand(msg tea.KeyPressMsg) (*model, tea.Cmd) {
	switch plumtypes.Key(msg.String()) {
	case plumtypes.CMD_CTRL_S:
		if m.state.Current() == CHAT_SECTION {
			return m, messages.Save(m.save)
		}
	}

	return m, nil
}

func (m *model) handleError(msg messages.Failed) (*model, tea.Cmd) {
	m.SetErr(msg.Err)
	return m, nil
}
