package onboardename

import (
	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/layout/header"
	"scne/app/shared/services/i19n"

	tea "charm.land/bubbletea/v2"
)

func (m model) View() tea.View {
	msg := i19n.Localize("view.not_found")
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
		m.form.View().Content,
	}
}

func (m *model) SetSectionSize(size size.Size) {
	// m.input.SetWidth(size.GetWidth())
}

func (m *model) GetSidepanel() []string {
	return []string{}
}
