package onboardename

import (
	tea "charm.land/bubbletea/v2"
)

// func (m *model) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
// 	switch msg.String() {
// 		case sharedtui.KEY_ENTER:
// 			name := strings.TrimSpace(m.name)
// 			if name == sharedtui.KEY_EMPTY {
// 				return m, nil
// 			}
// 			m.name = name

// 			return m, sharedtui.SendResolutionMsg(sharedtui.ResolutionMsg{
// 				Completed: true,
// 				Value:     m.name,
// 			})
// 		default:
// 			return m.handleDefaultInput(msg)
// 		}
// }

func (m *model) handleDefaultInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// newInput, cmd := m.input.Update(msg)
	// m.input = newInput
	// m.name = newInput.Value()
	return m, nil
}
func (m *model) delegateInputToForm(params tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := m.form.Update(params)
	return m, cmd
}
