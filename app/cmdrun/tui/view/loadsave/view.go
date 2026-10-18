package loadsave

import (
	"solopg/app/cmdrun/tui/shortcut"
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
		Title:    "onboarding:select.save.title",
		Subtitle: "onboarding:select.save.subtitle",
	}
}

func (m *model) GetSection() []string {
	return []string{
		m.list.View(),
	}
}

func (m *model) GetFooter() []string {
	return []string{
		shortcut.SELECT.String(),
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
