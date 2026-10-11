package gameboard

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/services/i19n"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)


func (m *model) GetSection() []string {

	base := m.viewport.View()

	if !m.isMenuActiveElementOpen() {
		base += "\n" + m.textarea.View()
	}
	view := tea.NewView(base)
	view.Cursor = getCursor(*m)
	view.AltScreen = true
	

	return []string{
		view.Content,
	}
}

func (m *model) SetSectionSize(size size.Size) {
}

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
