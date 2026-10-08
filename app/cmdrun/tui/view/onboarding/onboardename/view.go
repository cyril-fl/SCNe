package onboardename

import (
	tea "charm.land/bubbletea/v2"
)

func (m model) View() tea.View {
	// title := lipgloss.NewStyle().
	// 	Bold(true).
	// 	Render(i19n.Localize("onboarding"))

	// input := m.Input.View()

	// return tea.NewView(
	// 	lipgloss.JoinVertical(
	// 		lipgloss.Left,
	// 		/*
	// 			TODO LOW Verrifier pourquoi CREATE CHARACTER apparais pas sur les autre view
	// 		*/
	// 		title,
	// 		"",
	// 		transform.Capitalize(i19n.Localize("field.name")),
	// 		input,
	// 	),
	// )

	// page := page.NewPage(page.Template{
	// 	Title:    "onboarding:create_character",
	// 	Subtitle: "name",
	// 	Body:     m.input.View(),
	// })

	m.page.SetBody(m.input.View())

	return m.page.GetView()
}
