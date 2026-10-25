package combatboard

import (
	"solopg/app/cmdrun/tui/shortcut"
	"solopg/app/cmdrun/types/size"
)

// Getters & Setters

// GetView renders the combat chat.
func (m *model) GetView() string {
	return m.chat.GetView()
}

// GetFooter returns send and leave; the stat shortcuts are shown by the oracle section.
func (m *model) GetFooter() []string {
	return append(m.chat.GetFooter(), shortcut.LEAVE.String())
}

// SetSize sizes the combat chat and reloads the opponents; the gameboard calls it when the page is entered.
func (m *model) SetSize(s size.Size) {
	m.refresh()
	m.chat.SetSize(s)
}
