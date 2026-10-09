package section

import (
	"solopg/app/cmdrun/types/size"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Section interface {
	View() string
	Merge(model tea.Model) Section
	GetHeight() int
	GetWidth() int
	GetSize() size.Size
}

type section struct {
	parts   []string
	plugins []string
}

type HasSection interface {
	GetSection() []string
}

func New(parts ...string) Section {
	return &section{
		parts: parts,
	}
}

func (f *section) Merge(model tea.Model) Section {
	if view, ok := model.(HasSection); ok {
		f.plugins = append(f.plugins, view.GetSection()...)
	}
	return f
}

func (f *section) View() string {
	sectionparts := append(f.parts, f.plugins...)
	return strings.Join(sectionparts, " | ")
}

func (f *section) GetHeight() int {
	return lipgloss.Height(f.View())
}

func (f *section) GetWidth() int {
	return lipgloss.Width(f.View())
}

func (f *section) GetSize() size.Size {
	return size.New(size.Template{
		Height: f.GetHeight(),
		Width:  f.GetWidth(),
	})
}
