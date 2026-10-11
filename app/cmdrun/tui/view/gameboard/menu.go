package gameboard

import (
	"solopg/app/cmdrun/services/game"
	"solopg/app/cmdrun/tui/view/sidemenu/metadatamenu"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *model) GetSidepanel() []string {
	var sidemenu = []string{}
	for _, item := range m.menu {
		sidemenu = append(sidemenu, item.GetMenuView())
	}

	panel := makePanel(m.engine, sidemenu...)
	rendered := renderPanel(panel)

	return []string{
		rendered,
	}
}

func makePanel(engine *game.Engine, views ...string) string {
	statsview := makeStatView(engine)
	menuview := makeMenuView(views...)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		statsview,
		menuview,
	)
}

func makeStatView(engine *game.Engine) string {
	var content strings.Builder
	metadatamenu.GetCharacterInfo(&content, engine)
	metadatamenu.GetLocationInfo(&content, engine)
	metadatamenu.GetStatInfo(&content, engine)

	return strings.TrimRight(content.String(), "\n")
}

func makeMenuView(views ...string) string {
	menuview := []string{}
	for _, view := range views {
		menuview = append(menuview, view)
	}
	return strings.Join(menuview, "\n")
}

func renderPanel(content string) string {
	return lipgloss.NewStyle().
		Width(PANEL_WHIDTH-4).
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("5")).
		Render(content)
}
