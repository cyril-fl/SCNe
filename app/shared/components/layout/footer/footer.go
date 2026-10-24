package footer

import (
	"solopg/app/cmdrun/types/size"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Footer struct {
	parts   []string
	plugins []string
}

type HasFooter interface {
	GetFooter() []string
}

func New(parts ...string) *Footer {
	return &Footer{
		parts: parts,
	}
}

func Nil() *Footer {
	return &Footer{
		parts:   []string{},
		plugins: []string{},
	}
}

func (f *Footer) View() string {
	footerparts := append(f.parts, f.plugins...)
	return strings.Join(footerparts, " | ")
}

// Getters & Setters
func (f *Footer) GetHeight() int {
	return lipgloss.Height(f.View())
}

func (f *Footer) GetWidth() int {
	return lipgloss.Width(f.View())
}

func (f *Footer) GetSize() size.Size {
	return size.New(size.Template{
		Height: f.GetHeight(),
		Width:  f.GetWidth(),
	})
}

// Methods
func (f *Footer) Merge(model tea.Model) *Footer {
	if view, ok := model.(HasFooter); ok {
		f.plugins = view.GetFooter()
	}
	return f
}

// Helpers
