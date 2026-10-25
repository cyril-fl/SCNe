package cmdservetui

import (
	"errors"
	"os"

	"scne/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
)

// View renders the current step; step errors are logged silently.
func (m *model) View() tea.View {
	if m.steps.HasErr() {
		cwd, cwderr := os.Getwd()
		logs.SilentError("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, m.steps.GetErr()),
		})
	} else if current := m.steps.GetCurrentSubmodel(); current != nil {
		return current.View()
	}

	return tea.NewView("")
}
