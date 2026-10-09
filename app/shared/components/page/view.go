package page

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/services/i19n"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Pagineable interface {
	SetPageSize(size *tea.WindowSizeMsg)
}

type Page struct {
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

func NewPage(params Template) *Page {
	return &Page{
		title:    params.Title,
		subtitle: params.Subtitle,
		body:     params.Body,
	}
}

// Getters & Setters
func (p Page) GetView() tea.View {
	v := strings.Builder{}

	p.makeTitleGroup(&v)
	p.makeBody(&v)

	return tea.NewView(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Background(lipgloss.Color("2")).
			Render(v.String()),
	)
}

func (p *Page) SetSize(size *tea.WindowSizeMsg) {
	p.size = size
}

func (p *Page) SetBody(body string) {
	p.body = body
}

func (p *Page) GetSize() *tea.WindowSizeMsg {
	if p.size != nil {
		return p.size
	}
	return &tea.WindowSizeMsg{}
}

func (p *Page) GetAvailableSize() *tea.WindowSizeMsg {
	size := size.NewFromWindow(p.size)
	return &tea.WindowSizeMsg{
		Width:  size.GetWidth(),
		Height: max(0, size.GetHeight()-p.getTitleGroupHeight()),
	}
}

func (p Page) getTitleGroup() string {
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

func (p *Page) getTitleGroupHeight() int {
	if p.title == "" {
		return 0
	}
	return lipgloss.Height(p.getTitleGroup())
}

// Methods
func (p *Page) makeTitleGroup(v *strings.Builder) {
	if p.title != "" {
		v.WriteString(p.getTitleGroup())
		v.WriteString("\n")
	}
}

func (p *Page) makeBody(v *strings.Builder) {
	if p.body != "" {
		v.WriteString(p.body)
	}
}
