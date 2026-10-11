package gameboard

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/cmdrun/types/viewoptions"
	"solopg/app/shared/services/i19n"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
		"Section content",
	}
}

func (m *model) SetSectionSize(size size.Size) {
}




func (m model) GetFooter() []string {
	footer := []string{i19n.Localize("cmd.ctrl+s:save")}
	return append(footer, m.getMenuActiveElement().GetFooter()...)
}

// - Viewport - //
func (m *model) getContent() string {
	if m.isMenuActiveElementOpen() {
		return m.getMenuActiveElement().GetView()
	}

	welcomemsg := getWelcomeMessage(m)
	journalcontent := getJournalMessage(m)
	chatcontent := strings.Join(m.chat, "\n")

	composedView := lipgloss.JoinVertical(
		lipgloss.Left,

		lipgloss.NewStyle().
			Width(m.viewport.Width()).
			Foreground(lipgloss.Color("240")).
			Render(journalcontent.String()),

		lipgloss.NewStyle().
			Width(m.viewport.Width()).
			Render(welcomemsg.String()),

		lipgloss.NewStyle().
			Width(m.viewport.Width()).
			Render(chatcontent),
	)

	return composedView
}

func getWelcomeMessage(m *model) strings.Builder {
	welcomemsg := strings.Builder{}

	if len(m.chat) == 0 {
		welcomemsg.WriteString(i19n.Localize("chat.msg:welcome"))
		welcomemsg.WriteString("\n\n")
	}

	return welcomemsg
}

func getJournalMessage(m *model) strings.Builder {
	journalmsg := strings.Builder{}

	for _, entry := range m.adventurejournal {
		journalmsg.WriteString(entry)
		journalmsg.WriteString("\n")
	}

	if len(m.adventurejournal) > 0 {
		journalmsg.WriteString(strings.Repeat("_", m.viewport.Width()))
		journalmsg.WriteString("\n")
	}

	return journalmsg
}

// - Helper - //
// Refresh
func refreshViewport(m *model, params viewoptions.RefreshOption) (*model, tea.Cmd) {
	content := m.getContent()

	m.viewport.SetContent(content)

	if params.Positionreset {
		m.viewport.GotoTop()
	} else if params.ScrollBottom {
		m.viewport.GotoBottom()
	}

	return m, nil
}

func refreshLayout(m *model, msg tea.WindowSizeMsg) (*model, tea.Cmd) {
	chatWidth := max(0, msg.Width-PANEL_WHIDTH-PANEL_GAP)

	m.viewport.SetHeight(max(0, msg.Height-m.textarea.Height()-1))
	m.viewport.SetWidth(chatWidth)
	m.textarea.SetWidth(chatWidth)

	if !m.isMenuActiveElementOpen() {
		m.viewport.SetContent(m.getContent())
	}

	m.viewport.GotoBottom()

	return m, nil
}

// Render
func getCursor(m model) *tea.Cursor {
	cursor := m.textarea.Cursor()
	cursor.Y += lipgloss.Height(m.viewport.View())

	if item := m.getMenuActiveElement(); item != nil && item.IsOpen() {
		cursor = nil
	}

	return cursor
}
