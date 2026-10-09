package page

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	interfass "solopg/app/shared/types/interface"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type page struct {
	title    string
	subtitle string
	body     string
	cache
}

type cache struct {
	size *tea.WindowSizeMsg
}

type Template struct {
	Title    string
	Subtitle string
	Body     string
}

func New(params Template) interfass.Page {
	return &page{
		title:    params.Title,
		subtitle: params.Subtitle,
		body:     params.Body,
	}
}

// Getters & Setters
func (p page) GetView() tea.View {
	v := strings.Builder{}

	p.makeTitleGroup(&v)
	p.makeBody(&v)

	return tea.NewView(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			// Background(lipgloss.Color("2")).
			Render(p.handleFlex(&v)),
	)
}

func (p *page) SetBody(body string) {
	p.body = body
}

func (p *page) SetFooter(footer string) {
	logs.SilentWarning("error:not_implemented", map[string]any{
		"Function": "SetFooter",
		"Subject":  "page",
	})
}

func (p *page) GetSize() size.Size {
	if p.size != nil {
		return size.NewFromWindow(p.size)
	}
	return size.Nil()
}

func (p *page) SetSize(size *tea.WindowSizeMsg) {
	p.size = size
}

func (p *page) GetAvailableSize() size.Size {
	if p.size == nil {
		return size.Nil()
	}
	return size.New(size.Template{
		Width:  p.size.Width,
		Height: max(0, p.size.Height-p.getTitleGroupHeight()),
	})
}

func (p page) getTitleGroup() string {
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

func (p *page) getTitleGroupHeight() int {
	if p.title == "" {
		return 0
	}
	return lipgloss.Height(p.getTitleGroup())
}

// Methods
func (p *page) makeTitleGroup(v *strings.Builder) {
	if p.title != "" {
		v.WriteString(p.getTitleGroup())
		v.WriteString("\n")
	}
}

func (p *page) makeBody(v *strings.Builder) {
	if p.body != "" {
		v.WriteString(p.body)
	}
}

func (p *page) handleFlex(v *strings.Builder) string {
	size := p.GetSize()
	w := size.GetWidth()
	if w <= 0 {
		return v.String()
	}

	return lipgloss.NewStyle().
		Width(w).
		// Background(lipgloss.Color("2")).
		Render(v.String())
}
