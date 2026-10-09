package sidepanel


import (
	"solopg/app/cmdrun/types/size"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Sidepanel interface {
	View() string
	Merge(model tea.Model) Sidepanel
	GetHeight() int
	GetWidth() int
	GetSize() size.Size
}

type sidepanel struct {
	parts   []string
	plugins []string
}

type HasSidepanel interface {
	GetSidepanel() []string
}

func New(parts ...string) Sidepanel {
	return &sidepanel{
		parts: parts,
	}
}

func (f *sidepanel) Merge(model tea.Model) Sidepanel {
	if view, ok := model.(HasSidepanel); ok {
		f.plugins = append(f.plugins, view.GetSidepanel()...)
	}
	return f
}

func (f *sidepanel) View() string {
	footerparts := append(f.parts, f.plugins...)
	return strings.Join(footerparts, " | ")
}

func (f *sidepanel) GetHeight() int {
	return lipgloss.Height(f.View())
}

func (f *sidepanel) GetWidth() int {
	return lipgloss.Width(f.View())
}

func (f *sidepanel) GetSize() size.Size {
	return size.New(size.Template{
		Height: f.GetHeight(),
		Width:  f.GetWidth(),
	})
}
