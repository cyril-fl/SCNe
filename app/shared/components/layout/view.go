package layout

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/footer"
	"solopg/app/shared/components/page"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Layout struct {
	size   size.Size
	body   tea.Model
	footer footer.Footer
}

func New() Layout {
	return Layout{}
}

func NewWithSize(size size.Size) Layout {
	return Layout{
		size: size,
	}
}

// Getters & Setters
func (v *Layout) GetView() tea.View {
	// TODO le bord rose c'est la main view
	// Peu etre faire un factory avec ca du genre color, border, etc
	return tea.NewView(
		lipgloss.NewStyle().
			Height(v.getHeightWithinBorder()).
			Width(v.size.GetWidth()).
			// Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("5")).
			Render(v.makeView()),
	)
}

func (v *Layout) SetSize(size size.Size) {
	v.size = size
}

func (v *Layout) SetBody(content tea.Model) {
	v.body = content
	v.setBodyHeight(content)
}

func (v *Layout) setBodyHeight(model tea.Model) {
	if view, ok := model.(page.Pagineable); ok {
		b := trimPadding(v.body.View().Content)
		bodyBaseHeight := max(0, lipgloss.Height(b))
		bodyBaseHeight = 0
		view.SetPageSize(&tea.WindowSizeMsg{
			Width:  v.size.GetWidth(),
			Height: max(0, v.getHeightWithinBorder()-v.getFooterHight()-bodyBaseHeight),
		})
	}
}

func (v *Layout) SetFooter(parts ...string) {
	v.footer = footer.New(parts...)
}

func (v *Layout) getHeightWithinBorder() int {
	// offsetBecauseBorder := 2
	offsetBecauseBorder := 0

	return max(0, v.size.GetHeight()-offsetBecauseBorder)
}

func (v *Layout) getHeightAvailable() int {
	occupied := v.getFooterHight() + v.getBodyHeight()
	return max(0, v.getHeightWithinBorder()-occupied)
}

func (v *Layout) getBodyHeight() int {
	if v.body != nil {
		b := trimPadding(v.body.View().Content)
		return max(0, lipgloss.Height(b))
	}
	return 0
}

func (v *Layout) getFooterHight() int {
	if v.footer != nil {
		return v.footer.GetHeight()
	}
	return 0
}

// Methods
func (v *Layout) MergeFooter(other tea.Model) {
	if v.footer != nil {
		v.footer.Merge(other)
	}
}

func (v *Layout) makeView() string {
	view := strings.Builder{}
	v.makeBodyView(&view)
	v.makeFooterView(&view)
	return view.String()
}

func (v *Layout) makeBodyView(view *strings.Builder) {
	if v.body != nil {
		view.WriteString(trimPadding(v.body.View().Content))
	}
	// NOTE Explicit boody \n footer
	if v.footer != nil {
		// view.WriteString(fmt.Sprintf("\navH : %d | bH: %d. ", v.getHeightAvailable(), v.getBodyHeight()))
		view.WriteString("\n")
	}
}

func (v *Layout) makeFooterView(view *strings.Builder) {
	if v.footer != nil {
		view.WriteString(v.spaceFooter())
		view.WriteString(trimPadding(v.footer.View()))
	}
}

func (v *Layout) spaceFooter() string {
	maxHeight := v.getHeightAvailable()
	return strings.Repeat("\n", max(0, maxHeight))
}

// Helpers
func trimPadding(content string) string {
	return strings.TrimRight(content, "\n")
}
