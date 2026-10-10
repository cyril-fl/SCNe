package sidepanel

import (
	"solopg/app/cmdrun/types/size"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)
type Sidepanel struct {
	parts   []string
	plugins []string
}

type HasSidepanel interface {
	GetSidepanel() []string
}

func New(parts ...string) *Sidepanel {
	return &Sidepanel{
		parts: parts,
	}
}

func Nil() *Sidepanel {
	return &Sidepanel{
		parts:   []string{},
		plugins: []string{},
	}
}

func (f *Sidepanel) Merge(model tea.Model) *Sidepanel {
	if view, ok := model.(HasSidepanel); ok {
		f.plugins = append(f.plugins, view.GetSidepanel()...)
	}
	return f
}

func (f *Sidepanel) View() string {
	footerparts := append(f.parts, f.plugins...)
	return strings.Join(footerparts, " | ")
}

func (f *Sidepanel) GetHeight() int {
	return lipgloss.Height(f.View())
}

func (f *Sidepanel) GetWidth() int {
	return lipgloss.Width(f.View())
}

func (f *Sidepanel) GetSize() size.Size {
	return size.New(size.Template{
		Height: f.GetHeight(),
		Width:  f.GetWidth(),
	})
}
