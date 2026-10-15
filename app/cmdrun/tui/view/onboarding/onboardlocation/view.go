package onboardlocation

import (
	"fmt"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/layout/header"
	"solopg/app/shared/services/i19n"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/utils/transform"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	msg := "404 - Not Found"
	return tea.NewView(msg)
}

func (m *model) GetHeader() header.Template {
	return header.Template{
		Title:    "onboarding.character:create",
		Subtitle: "onboarding:select.location",
	}
}

func (m *model) GetSection() []string {
	return []string{
		m.makeBody(),
	}
}

func (m *model) SetSectionSize(size size.Size) {
	informationHeight := lipgloss.Height(m.makeInformationView())
	separatorHeight := lipgloss.Height("\n")
	listHeight := max(0, size.GetHeight()-informationHeight-separatorHeight)

	m.reroll.Choice.SetSize(size.GetWidth(), listHeight)
}

func (m *model) GetSidepanel() []string {
	return []string{
		"Sidepanel content",
	}
}

func (m *model) GetFooter() []string {
	return []string{
		"Enter: select",
	}
}

func (m model) makeBody() string {
	v := lipgloss.JoinVertical(
		lipgloss.Left,
		m.makeInformationView(),
		"\n",
		m.reroll.Choice.View(),
	)

	return lipgloss.NewStyle().
		// Background(lipgloss.Color("2")).
		Render(v)
}

func (m *model) makeInformationView() string {
	v := strings.Builder{}
	v.WriteString(fmt.Sprintf("%s : ", i19n.Localize(m.reroll.Value.GetName())))

	if sharedtui.EventKey(m.reroll.Value.Descriptible.GetDescription()) != sharedtui.KEY_EMPTY {
		v.WriteString(i19n.Localize(m.reroll.Value.Descriptible.GetDescription()))
	}

	v.WriteString(strings.Join([]string{"| ", transform.Capitalize(i19n.Localize("attempt")), ": "}, ""))
	v.WriteString(fmt.Sprintf("%d", m.reroll.Attempt))

	return v.String()
}
