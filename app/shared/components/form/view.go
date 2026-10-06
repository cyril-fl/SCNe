package form

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Form) View() tea.View {
	contents := make([]string, 0, len(m.fields))

	for _, f := range m.fields {
		contents = append(contents, f.View().Content)
	}

	if m.GetError() != nil {
		err := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Render("\n" + m.GetError().Error())

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

		if m.GetError() != nil {
			err := lipgloss.NewStyle().
				Width(viewport.Width()).
				Foreground(lipgloss.Color("1")).
				Render("\n" + m.GetError().Error())

			height += lipgloss.Height(err)
		}

		return top, height

	}

	return -1, 0
}
