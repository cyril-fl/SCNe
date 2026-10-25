package sidepanel

import (
	"solopg/app/cmdrun/types/size"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Sidepanel renders the side area: its own parts, then those of the merged model.
type Sidepanel struct {
	parts   []string
	plugins []string
}

// HasSidepanel is implemented by models that render in the side panel.
type HasSidepanel interface {
	GetSidepanel() []string
}

// HasPanelSize is implemented by models that need the side panel size.
type HasPanelSize interface {
	SetPanelSize(size size.Size)
}

// New returns a side panel with fixed parts.
func New(parts ...string) *Sidepanel {
	return &Sidepanel{
		parts: parts,
	}
}

// Nil returns an empty side panel.
func Nil() *Sidepanel {
	return &Sidepanel{
		parts:   []string{},
		plugins: []string{},
	}
}

// Merge replaces the plugin parts with the model lines, joined as one part.
func (s *Sidepanel) Merge(model tea.Model) *Sidepanel {
	if view, ok := model.(HasSidepanel); ok {
		s.plugins = []string{strings.Join(view.GetSidepanel(), "\n")}
	}
	return s
}

// View joins the parts with " | ".
func (s *Sidepanel) View() string {
	footerparts := append(s.parts, s.plugins...)
	return strings.Join(footerparts, " | ")
}

// GetHeight returns the rendered height.
func (s *Sidepanel) GetHeight() int {
	return lipgloss.Height(s.View())
}

// GetWidth returns the rendered width.
func (s *Sidepanel) GetWidth() int {
	return lipgloss.Width(s.View())
}

// GetSize returns the rendered size.
func (s *Sidepanel) GetSize() size.Size {
	return size.New(size.Template{
		Height: s.GetHeight(),
		Width:  s.GetWidth(),
	})
}

// SetSize forwards size to model when it implements HasPanelSize.
func (s *Sidepanel) SetSize(model tea.Model, size size.Size) {
	if view, ok := model.(HasPanelSize); ok {
		view.SetPanelSize(size)
	}
}
