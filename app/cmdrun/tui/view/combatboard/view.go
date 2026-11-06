package combatboard

import (
	"scne/app/cmdrun/tui/shortcut"
	"scne/app/cmdrun/types/size"
)

// Getters & Setters

// GetView renders the active plugin, or the combat chat.
func (m *model) GetView() string {
	if active := m.activePlugin(); active != nil {
		return active.GetView()
	}
	return m.chat.GetView()
}

// GetFooter returns the active plugin shortcuts, or the chat, leave and plugin
// shortcuts; the stat shortcuts are shown by the oracle section.
func (m *model) GetFooter() []string {
	if active := m.activePlugin(); active != nil {
		return active.GetFooter()
	}

	footer := append(m.chat.GetFooter(), shortcut.LEAVE.String())
	for _, p := range m.plugins {
		footer = append(footer, p.GetFooter()...)
	}
	return footer
}

// SetSize sizes the combat chat and the plugins and reloads the opponents; the gameboard calls it when the page is entered.
func (m *model) SetSize(s size.Size) {
	m.refresh()
	m.chat.SetSize(s)
	for _, p := range m.plugins {
		p.SetSize(s)
	}
}
