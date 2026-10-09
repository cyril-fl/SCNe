package onboardlocation

import (
	"fmt"
	"solopg/app/shared/services/i19n"
	sharedtui "solopg/app/shared/tui"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// TODO Ui ici bug !
func (m model) View() tea.View {
	m.page.SetBody(m.makeBody())
	return m.page.GetView()
}

func (m model) makeBody() string {
	v := strings.Builder{}

	m.makeInformationView(&v)
	m.makeOptionsView(&v)

	return lipgloss.NewStyle().
		// Background(lipgloss.Color("2")).
		Render(v.String())
}

func (m *model) makeInformationView(v *strings.Builder) {
	v.WriteString(strings.Join([]string{"", i19n.Localize("attempt"), ": "}, ""))
	v.WriteString(fmt.Sprintf("%d\n", m.reroll.Attempt))

	v.WriteString(fmt.Sprintf("%s : ", i19n.Localize(m.reroll.Value.Name)))
	if m.reroll.Value.Description.Description != sharedtui.KEY_EMPTY {
		v.WriteString(i19n.Localize(m.reroll.Value.Description.Description))
	}
}

func (m *model) makeOptionsView(v *strings.Builder) {
	v.WriteString("\n\n")
	v.WriteString(m.reroll.Options.View())
}

func (m *model) SetPageSize(size *tea.WindowSizeMsg) {
	m.page.SetSize(size)

	bodySize := m.page.GetAvailableSize()
	// TODO try to improve
	offsetInfoViewSize := 2
	offetAlinea := 1
	offset := offsetInfoViewSize + offetAlinea

	m.reroll.Options.SetSize(bodySize.Width, bodySize.Height- offset)
}


