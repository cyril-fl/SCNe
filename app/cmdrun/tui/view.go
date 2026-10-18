package cmdruntui

import (
	"errors"
	"os"
	"solopg/app/cmdrun/tui/shortcut"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
)

func (m model) View() tea.View {
	m.layout.SetSize(size.NewFromWindow(m.steps.GetSize()))
	m.layout.SetFooter(shortcut.QUIT.String())

	if m.steps.HasErr() {
		cwd, cwderr := os.Getwd()
		logs.SilentError("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, m.steps.GetErr()),
		})
	} else if current := m.steps.GetCurrentSubmodel(); current != nil {
		m.layout.SetHeader(current)
		m.layout.MergeFooter(current)
		m.layout.MergeSection(current)
		m.layout.MergeSidepanel(current)
	}

	return m.layout.GetView()
}
