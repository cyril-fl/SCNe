package onboardarchetype

import (
	tea "charm.land/bubbletea/v2"
)

func (m model) View() tea.View {
	bodySize := m.page.GetAvailableSize()
	m.list.SetSize(bodySize.Width, bodySize.Height)

	m.page.SetBody(m.list.View())

	return m.page.GetView()
}

func (m model) SetPageSize(size *tea.WindowSizeMsg) {
	m.page.SetSize(size)
}
