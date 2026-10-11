package cmdruntui

import (
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"

	tea "charm.land/bubbletea/v2"
)

// Handlers
func (m *model) handleEvent(msg tea.KeyMsg) (*model, tea.Cmd) {
	switch sharedtui.EventKey(msg.String()) {
	case sharedtui.SHORT_CTRL_Q:
		logs.SilentInfo("system.msg:quit")
		return m, tea.Quit
	case sharedtui.KEY_ESC:
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
