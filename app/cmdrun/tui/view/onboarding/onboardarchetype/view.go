package onboardarchetype

import (
	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/layout/header"

	tea "charm.land/bubbletea/v2"
)

func (m model) View() tea.View {
	msg := "404 - Not Found"
	return tea.NewView(msg)
}

func (m *model) GetHeader() header.Template {
	return header.Template{
		Title:    "onboarding.character:create",
		Subtitle: "onboarding:select." + m.name,
	}
}

func (m *model) GetSection() []string {
	return []string{
		m.list.View(),
	}
}

func (m *model) SetSectionSize(size size.Size) {
	m.list.SetSize(size.GetWidth(), size.GetHeight())
}

func (m *model) GetSidepanel() []string {
	return []string{
		"Sidepanel content",
	}
}
