package cmdruntui

import (
	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
	"scne/app/shared/services/logs"
)

// Handlers
func (m *model) handleEvent(msg tea.KeyMsg) (*model, tea.Cmd) {
	switch plumtypes.Key(msg.String()) {
	case plumtypes.SHORT_CTRL_Q:
		logs.SilentInfo("system.msg:quit")
		return m, tea.Quit
	case plumtypes.KEY_ESC:
		current := m.steps.GetCurrentSubmodel()
		if current == nil {
			logs.SilentInfo("system.msg:quit")
			return m, tea.Quit
		}
		if current, ok := current.(EscapeSupport); !ok || (ok && !current.HandlesEscape()) {
			logs.SilentInfo("system.msg:quit")
			return m, tea.Quit
		}
	}
	return m, m.steps.UpdateSubmodel(msg)
}
