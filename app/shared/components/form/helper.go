package form

import (
	sharedtui "solopg/app/shared/tui"

	tea "charm.land/bubbletea/v2"
)

func (m *Form) handleKeyPressMsg(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch sharedtui.EventKey(msg.String()) {
	case sharedtui.KEY_TAB:
		return m, Send[NextField]()
	case sharedtui.KEY_SHIFT_TAB:
		return m, Send[PreviousField]()
	case m.sendKey:
		return m, Send[Validate]()
	}
	return m.updateFocusedField(msg)
}

func (m *Form) handleErrorsReset() {
	m.ClearErr()
	for _, f := range m.fields {
		f.ResetError()
	}
}

func (m *Form) updateFocusedField(msg tea.Msg) (tea.Model, tea.Cmd) {
	index := m.getCurrentFieldIndex()
	if index < 0 {
		return m, nil
	}

	newField, cmd := m.fields[index].Update(msg)
	m.fields[index] = newField.(Field)

	return m, cmd
}

func (m *Form) focusOnIndex(index int) *Form {
	for _, f := range m.fields {
		f.Blur()
	}
	if index >= 0 && index < len(m.fields) {
		m.fields[index].Focus()
	} else if len(m.fields) > 0 {
		m.fields[0].Focus()
	}

	return m
}

func (m *Form) isLastField() bool {
	currentIndex := m.getCurrentFieldIndex()
	return currentIndex == len(m.fields)-1
}
