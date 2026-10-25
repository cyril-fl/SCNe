package logstream

import (
	"scne/app/shared/services/factory/millfilterlog"
	"scne/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
)

// View renders the received logs filtered by --filter.
func (m *model[T]) View() tea.View {
	mill := millfilterlog.New(millfilterlog.Template{
		LogList: logs.ConvertStringablesToList(m.stream.GetEvents()),
		Flags:   m.flags,
	})

	mill.Filter()

	return tea.NewView(mill.GetJoinedStringList())
}
