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

type HasPanelSize interface {
	SetPanelSize(size size.Size)
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

func (s *Sidepanel) Merge(model tea.Model) *Sidepanel {
	if view, ok := model.(HasSidepanel); ok {
		s.plugins = []string{strings.Join(view.GetSidepanel(), "\n")}
	}
	return s
}

func (s *Sidepanel) View() string {
	footerparts := append(s.parts, s.plugins...)
	return strings.Join(footerparts, " | ")
}

func (s *Sidepanel) GetHeight() int {
	return lipgloss.Height(s.View())
}

func (s *Sidepanel) GetWidth() int {
	return lipgloss.Width(s.View())
}

func (s *Sidepanel) GetSize() size.Size {
	return size.New(size.Template{
		Height: s.GetHeight(),
		Width:  s.GetWidth(),
	})
}

func (s *Sidepanel) SetSize(model tea.Model, size size.Size) {
	if view, ok := model.(HasPanelSize); ok {
		view.SetPanelSize(size)
	}
}
