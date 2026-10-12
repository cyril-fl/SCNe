package gameboard

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/services/i19n"

	tea "charm.land/bubbletea/v2"
)

/*
TODO LOW Trouver un meilleur moyen de gerer les sizes,
trouver un moyen de les distibuer al'interieur des sous composant avec
par exemple un contexte disponible, un peu comme Engine et States.
*/
const (
	PANEL_WHIDTH  = 34
	PANEL_GAP     = 1
	ORACLE_HEIGHT = 1
)

func (m model) View() tea.View {
	msg := "404 - Not Found"
	return tea.NewView(msg)
}

func (m *model) GetSection() []string {
	return []string{
		m.state.GetView(),
	}
}

func (m *model) SetSectionSize(size size.Size) {
	m.state.SetSize(size)
}

func (m model) GetFooter() []string {
	footer := []string{i19n.Localize("cmd.ctrl+s:save")}
	return append(footer, m.state.getFooter()...)
}
