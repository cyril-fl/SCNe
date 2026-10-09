package layout

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/layout/footer"
	"solopg/app/shared/components/layout/header"
	"solopg/app/shared/components/layout/section"
	"solopg/app/shared/components/layout/sidepanel"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Layout struct {
	size   size.Size
	// body   tea.Model
	header header.Header
	sidepanel sidepanel.Sidepanel
	section section.Section
	footer footer.Footer
}

func New() Layout {
	return Layout{}
}

// func NewWithSize(size size.Size) Layout {
// 	return Layout{
// 		size: size,
// 	}
// }

// Getters & Setters
func (v *Layout) GetView() tea.View {
	return tea.NewView(
		lipgloss.NewStyle().
			Height(v.getHeightWithinBorder()).
			Width(v.getWidthtWithinBorder()).
			// TODO le bord rose c'est la main view
			// Peu etre faire un factory avec ca du genre color, border, etc
			// Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("5")).
			Render(v.makeView()),
	)
}

func (v *Layout) GetSize() size.Size {
	return v.size
}

func (v *Layout) SetSize(size size.Size) {
	v.size = size
}

func (v *Layout) GetAvailableSize() size.Size {
	return size.New(size.Template{
		Width:  v.size.GetWidth(),
		Height: v.getHeightAvailable(),
	})
}

func (v *Layout) getHeightAvailable() int {
	occupied := v.getFooterHight() + v.getSectionHeight()
	return max(0, v.getHeightWithinBorder()-occupied)
}

func (v *Layout) getWidthAvailable() int {
	// occupied := v.getFooterHight() + v.getBodyHeight()
	return max(0, v.getWidthtWithinBorder())
}

func (v *Layout) getHeightWithinBorder() int {
	// offsetBecauseBorder := 2
	offsetBecauseBorder := 0
	return max(0, v.size.GetHeight()-offsetBecauseBorder)
}

func (v *Layout) getWidthtWithinBorder() int {
	// offsetBecauseBorder := 2
	offsetBecauseBorder := 0
	return max(0, v.size.GetWidth()-offsetBecauseBorder)
}

// Methods
func (v *Layout) makeView() string {
	view := strings.Builder{}
	v.makeHeaderView(&view)
	v.makeSidepanelView(&view)
	v.makeSectionView(&view)
	v.makeFooterView(&view)
	return view.String()
}

// Helpers
func trimParagraphPadding(content string) string {
	return strings.TrimRight(content, "\n")
}
