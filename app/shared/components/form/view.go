package form

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"solopg/app/shared/services/style"
)

func (m *Form) View() tea.View {
	contents := make([]string, 0, len(m.fields))

	for _, f := range m.fields {
		contents = append(contents, f.View().Content)
	}

	if m.HasErr() {
		err := style.Error.Render("\n" + m.GetErr().Error())

		contents = append(contents, "\n"+err)
	}

	return tea.NewView(
		lipgloss.JoinVertical(lipgloss.Left, contents...),
	)
}

func (m *Form) GetFieldScrollArea(viewport viewport.Model) (int, int) {
	top := 0

	for _, field := range m.fields {
		content := lipgloss.NewStyle().Width(viewport.Width()).Render(field.View().Content)
		height := lipgloss.Height(content)

		if !field.IsFocused() {
			top += height
			continue
		}

		if m.HasErr() {
			err := style.Error.Width(viewport.Width()).Render("\n" + m.GetErr().Error())

			height += lipgloss.Height(err)
		}

		return top, height

	}

	return -1, 0
}
