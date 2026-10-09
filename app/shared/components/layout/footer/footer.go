package footer

import (
	"solopg/app/cmdrun/types/size"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Footer interface {
	View() string
	Merge(model tea.Model) Footer
	GetHeight() int
	GetWidth() int
	GetSize() size.Size
}

type footer struct {
	parts   []string
	plugins []string
}

type HasFooter interface {
	GetFooter() []string
}

func New(parts ...string) Footer {
	return &footer{
		parts: parts,
	}
}

func (f *footer) Merge(model tea.Model) Footer {
	if view, ok := model.(HasFooter); ok {
		f.plugins = append(f.plugins, view.GetFooter()...)
	}
	return f
}

func (f *footer) View() string {
	footerparts := append(f.parts, f.plugins...)
	return strings.Join(footerparts, " | ")
}

func (f *footer) GetHeight() int {
	return lipgloss.Height(f.View())
}

func (f *footer) GetWidth() int {
	return lipgloss.Width(f.View())
}

func (f *footer) GetSize() size.Size {
	return size.New(size.Template{
		Height: f.GetHeight(),
		Width:  f.GetWidth(),
	})
}
