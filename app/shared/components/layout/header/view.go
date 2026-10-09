package header


import (
	"solopg/app/cmdrun/types/size"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Header interface {
	View() string
	Merge(model tea.Model) Header
	GetHeight() int
	GetWidth() int
	GetSize() size.Size
}

type header struct {
	parts   []string
	plugins []string
}

type HasHeader interface {
	GetHeader() []string
}
	
func New(parts ...string) Header {
	return &header{
		parts: parts,
	}
}

func (f *header) Merge(model tea.Model) Header {
	if view, ok := model.(HasHeader); ok {
		f.plugins = append(f.plugins, view.GetHeader()...)
	}
	return f
}

func (f *header) View() string {
	footerparts := append(f.parts, f.plugins...)
	return strings.Join(footerparts, " | ")
}

func (f *header) GetHeight() int {
	return lipgloss.Height(f.View())
}

func (f *header) GetWidth() int {
	return lipgloss.Width(f.View())
}

func (f *header) GetSize() size.Size {
	return size.New(size.Template{
		Height: f.GetHeight(),
		Width:  f.GetWidth(),
	})
}
