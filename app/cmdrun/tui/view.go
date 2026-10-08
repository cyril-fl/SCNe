package cmdruntui

import (
	"errors"
	"os"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
)

func (m model) View() tea.View {
	m.mainview.SetSize(size.NewFromWindow(m.size))
	m.mainview.SetFooter(i19n.Localize("cmd.ctrl+q:quit"))

	if m.err != nil {
		cwd, cwderr := os.Getwd()
		_ = logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, m.err),
		})

		// m.mainview.SetBody(newErrorModel(err))
	} else if current := m.steps.GetCurrentSubmodel(); current != nil {
		m.mainview.MergeFooter(current)
		m.mainview.SetBody(current)
	}

	return m.mainview.GetView()
}
