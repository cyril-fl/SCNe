package codexmenu

import (
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/utils/transform"

	"charm.land/lipgloss/v2"
)

func (m *codexMenu) GetMenuView() string {
	title := transform.Uppercase(i19n.Localize(m.ID()))
	title = lipgloss.NewStyle().Bold(true).Render(title)
	return lipgloss.JoinVertical(lipgloss.Left, title, "", m.GetList().View())
}

// Content is rendered by the gameboard state machine in the main section.
func (m *codexMenu) GetView() string { return "" }

func (m *codexMenu) GetFooter() []string {
	return []string{i19n.Localize("cmd.alt+enter:open")}
}
