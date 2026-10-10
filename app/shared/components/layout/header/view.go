package header

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/services/i19n"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Header struct {
	title    string
	subtitle string
}

type Template struct {
	Title    string
	Subtitle string
}

type HasHeader interface {
	GetHeader() Template
}

func New(params Template) *Header {
	return &Header{
		title:    params.Title,
		subtitle: params.Subtitle,
	}
}

func Nil() *Header {
	return &Header{
		title:    "",
		subtitle: "",
	}
}

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
func (p *Header) SetHeader(model tea.Model) {
	if view, ok := model.(HasHeader); ok {
		template := view.GetHeader()

		p.title = template.Title
		p.subtitle = template.Subtitle
	}
}

func (p *Header) GetHeight() int {
	return lipgloss.Height(p.View())
}

func (p *Header) GetWidth() int {
	return lipgloss.Width(p.View())
}

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

	return lipgloss.NewStyle().
		Bold(true).
		Render(view)
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