package onboardename

import tea "charm.land/bubbletea/v2"

func (m *model) handleDefaultInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	newInput, cmd := m.input.Update(msg)
	m.input = newInput
	m.name = newInput.Value()
	return m, cmd
}