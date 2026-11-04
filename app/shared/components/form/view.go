package form

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"scne/app/shared/services/style"
)

// View renders the fields stacked vertically, followed by the form error if any.
func (m *Form) View() tea.View {
	contents := make([]string, 0, len(m.fields))

	for _, f := range m.fields {
		contents = append(contents, f.View().Content)
	}

	if m.HasErr() {
		err := style.Error().Render("\n" + m.GetErr().Error())

		contents = append(contents, "\n"+err)
	}

	return tea.NewView(
		m.scrollToFocus(lipgloss.JoinVertical(lipgloss.Left, contents...)),
	)
}

// scrollToFocus returns content cut to the form height, moved only as far as needed to show the focused field.
func (m *Form) scrollToFocus(content string) string {
	if m.height <= 0 || lipgloss.Height(content) <= m.height {
		return content
	}

	m.view.SetWidth(lipgloss.Width(content))
	m.view.SetHeight(m.height)
	m.view.SetContent(content)

	if top, height := m.GetFieldScrollArea(m.view); top >= 0 {
		// The top wins when the field is taller than the form.
		m.view.SetYOffset(min(top, max(m.view.YOffset(), top+height-m.height)))
	}

	return m.view.View()
}

// GetFieldScrollArea returns the top offset and height of the focused field in viewport.
// It returns -1, 0 when no field is focused.
func (m *Form) GetFieldScrollArea(viewport viewport.Model) (int, int) {
	top := 0

	for _, field := range m.fields {
		content := lipgloss.NewStyle().Width(viewport.Width()).Render(field.View().Content)
		height := lipgloss.Height(content)

		if !field.IsFocused() {
			top += height
			continue
		}

		if repeat, ok := field.(*repeatField); ok {
			innerTop, innerHeight := repeat.group.GetFieldScrollArea(viewport)
			return top + innerTop, innerHeight
		}

		if m.HasErr() {
			err := style.Error().Width(viewport.Width()).Render("\n" + m.GetErr().Error())

			height += lipgloss.Height(err)
		}

		return top, height

	}

	return -1, 0
}
