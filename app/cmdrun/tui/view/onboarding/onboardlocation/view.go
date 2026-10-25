package onboardlocation

import (
	"fmt"
	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/layout/header"
	"scne/app/shared/services/i19n"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	plumtypes "github.com/cyril-fl/plum/types"
	plumtransform "github.com/cyril-fl/plum/utils/transform"
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

	if plumtypes.Key(m.reroll.Value.Descriptible.GetDescription()) != plumtypes.KEY_EMPTY {
		v.WriteString(i19n.Localize(m.reroll.Value.Descriptible.GetDescription()))
	}

	v.WriteString(strings.Join([]string{"| ", plumtransform.Capitalize(i19n.Localize("attempt")), ": "}, ""))
	v.WriteString(fmt.Sprintf("%d", m.reroll.Attempt))

	return v.String()
}
