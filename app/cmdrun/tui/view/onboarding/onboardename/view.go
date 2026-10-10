package onboardename

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/layout/header"

	tea "charm.land/bubbletea/v2"
)

func (m model) View() tea.View {
	msg := "404 - Not Found"
	return tea.NewView(msg)
}

func (m *model) GetHeader() header.Template {
	return header.Template{
		Title:    "onboarding.character:create",
		Subtitle: "onboarding:select.character.name",
	}
}

func (m *model) GetSection() []string {
	return []string{
		m.input.View(),
	}
}

func (m *model) SetSectionSize(size size.Size) {
	m.input.SetWidth(size.GetWidth())
}

func (m *model) GetFooter() []string {
	return []string{
		"TEST",
	}
}
