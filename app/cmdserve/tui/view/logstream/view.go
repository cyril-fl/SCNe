package logstream

import (
	"solopg/app/shared/services/factory/millfilterlog"
	"solopg/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
)

func (m *model[T]) View() tea.View {
	mill := millfilterlog.New(millfilterlog.Template{
		LogList: logs.ConvertStringablesToList(m.stream.GetEvents()),
		Flags:   m.flags,
	})

	mill.Filter()

	return tea.NewView(mill.GetJoinedStringList())
}
