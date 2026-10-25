package section

import (
	"solopg/app/cmdrun/types/size"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Section renders the main area: its own parts, then those of the merged model.
type Section struct {
	parts   []string
	plugins []string
}

// HasSection is implemented by models that render in the section.
type HasSection interface {
	GetSection() []string
}

// HasSectionSize is implemented by models that need the section size.
type HasSectionSize interface {
	SetSectionSize(size size.Size)
}

// New returns a section with fixed parts.
func New(parts ...string) *Section {
	return &Section{
		parts: parts,
	}
}

// Nil returns an empty section.
func Nil() *Section {
	return &Section{
		parts:   []string{},
		plugins: []string{},
	}
}

// View stacks the parts vertically.
func (s *Section) View() string {
	v := append(s.parts, s.plugins...)
	return lipgloss.JoinVertical(lipgloss.Left, v...)
}

// Getters & Setters

// GetHeight returns the rendered height.
func (s *Section) GetHeight() int {
	return lipgloss.Height(s.View())
}

// GetWidth returns the rendered width.
func (s *Section) GetWidth() int {
	return lipgloss.Width(s.View())
}

// GetSize returns the rendered size.
func (s *Section) GetSize() size.Size {
	return size.New(size.Template{
		Height: s.GetHeight(),
		Width:  s.GetWidth(),
	})
}

// Methods

// Merge replaces the plugin parts with the model section when it implements HasSection.
func (s *Section) Merge(model tea.Model) *Section {
	if view, ok := model.(HasSection); ok {
		s.plugins = view.GetSection()
	}
	return s
}

// SetSize forwards size to model when it implements HasSectionSize.
func (s *Section) SetSize(model tea.Model, size size.Size) {
	if view, ok := model.(HasSectionSize); ok {
		view.SetSectionSize(size)
	}
}

// Helpers
