package section

import (
	"solopg/app/cmdrun/types/size"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Section struct {
	parts   []string
	plugins []string
}

type HasSection interface {
	GetSection() []string
	SetSectionSize(size size.Size)
}

func New(parts ...string) *Section {
	return &Section{
		parts: parts,
	}
}

func Nil() *Section {
	return &Section{
		parts:   []string{},
		plugins: []string{},
	}
}

func (s *Section) View() string {
	v := append(s.parts, s.plugins...)
	return lipgloss.JoinVertical(lipgloss.Left, v...)
}

// Getters & Setters
func (s *Section) GetHeight() int {
	return lipgloss.Height(s.View())
}

func (s *Section) GetWidth() int {
	return lipgloss.Width(s.View())
}

func (s *Section) GetSize() size.Size {
	return size.New(size.Template{
		Height: s.GetHeight(),
		Width:  s.GetWidth(),
	})
}

// Methods
func (s *Section) Merge(model tea.Model) *Section {
	if view, ok := model.(HasSection); ok {
		s.plugins = append(s.plugins, view.GetSection()...)
	}
	return s
}

func (s *Section) SetSize(model tea.Model, size size.Size) {
	if view, ok := model.(HasSection); ok {
		view.SetSectionSize(size)
	}
}

// Helpers
