package footer

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/services/style"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var separator = style.Muted.Render(" · ")

// Footer renders a one-line shortcut bar: its own parts, then those of the merged model.
type Footer struct {
	parts   []string
	plugins []string
}

// HasFooter is implemented by models that add shortcuts to the footer.
type HasFooter interface {
	GetFooter() []string
}

// New returns a footer with fixed parts.
func New(parts ...string) *Footer {
	return &Footer{
		parts: parts,
	}
}

// Nil returns an empty footer.
func Nil() *Footer {
	return &Footer{
		parts:   []string{},
		plugins: []string{},
	}
}

// View joins the parts with a muted " · ".
func (f *Footer) View() string {
	footerparts := append(f.parts, f.plugins...)
	return strings.Join(footerparts, separator)
}

// Getters & Setters

// GetHeight returns the rendered height.
func (f *Footer) GetHeight() int {
	return lipgloss.Height(f.View())
}

// GetWidth returns the rendered width.
func (f *Footer) GetWidth() int {
	return lipgloss.Width(f.View())
}

// GetSize returns the rendered size.
func (f *Footer) GetSize() size.Size {
	return size.New(size.Template{
		Height: f.GetHeight(),
		Width:  f.GetWidth(),
	})
}

// Methods

// Merge replaces the plugin parts with the model footer when it implements HasFooter.
func (f *Footer) Merge(model tea.Model) *Footer {
	if view, ok := model.(HasFooter); ok {
		f.plugins = view.GetFooter()
	}
	return f
}

// Helpers
