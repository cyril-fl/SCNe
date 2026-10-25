package header

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/style"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Header renders a localized title and subtitle.
type Header struct {
	title    string
	subtitle string
}

// Template holds the header translation keys.
type Template struct {
	Title    string
	Subtitle string
}

// HasHeader is implemented by models that provide a header.
type HasHeader interface {
	GetHeader() Template
}

// New returns a header from params.
func New(params Template) *Header {
	return &Header{
		title:    params.Title,
		subtitle: params.Subtitle,
	}
}

// Nil returns an empty header.
func Nil() *Header {
	return &Header{
		title:    "",
		subtitle: "",
	}
}

// View renders the title group, or an empty string without a title.
func (p *Header) View() string {
	v := strings.Builder{}

	p.makeTitleGroup(&v)

	hgroup := tea.NewView(
		lipgloss.NewStyle().
			Render(p.handleFlex(&v)),
	)

	return hgroup.Content
}

// Getters & Setters

// SetHeader copies the model header when it implements HasHeader.
func (p *Header) SetHeader(model tea.Model) {
	if view, ok := model.(HasHeader); ok {
		template := view.GetHeader()

		p.title = template.Title
		p.subtitle = template.Subtitle
	}
}

// GetHeight returns the rendered height.
func (p *Header) GetHeight() int {
	return lipgloss.Height(p.View())
}

// GetWidth returns the rendered width.
func (p *Header) GetWidth() int {
	return lipgloss.Width(p.View())
}

// GetSize returns the rendered size.
func (p *Header) GetSize() size.Size {
	return size.New(size.Template{
		Height: p.GetHeight(),
		Width:  p.GetWidth(),
	})
}

func (p Header) getTitleGroup() string {
	if p.title == "" {
		return ""
	}

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		i19n.Localize(p.title),
		i19n.Localize(p.subtitle),
		"",
	)

	return style.Title.Render(view)
}

func (p *Header) getTitleGroupHeight() int {
	if p.title == "" {
		return 0
	}
	return lipgloss.Height(p.getTitleGroup())
}

// Methods
func (p *Header) makeTitleGroup(v *strings.Builder) {
	if p.title != "" {
		v.WriteString(p.getTitleGroup())
		v.WriteString("\n")
	}
}

func (p *Header) handleFlex(v *strings.Builder) string {
	w := lipgloss.Width(v.String())
	if w <= 0 {
		return v.String()
	}

	return lipgloss.NewStyle().
		Width(w).
		// Background(lipgloss.Color("2")).
		Render(v.String())
}

// Helpers
