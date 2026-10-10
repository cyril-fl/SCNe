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
	size      size.Size
	header    *header.Header
	sidepanel *sidepanel.Sidepanel
	section   *section.Section
	footer    *footer.Footer
}

const (
	// offsetBecauseBorder = 2
	offsetBecauseBorder = 0
)

func New() Layout {
	return Layout{}
}

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
	return max(0, v.getHeightWithinBorder()-v.getHeightOccupied())
}

func (v *Layout) getHeightOccupied() int {
	occupied := []int{
		v.getFooterHeight(),
		v.getSectionHeight(),
		v.getHeaderHeight(),
	}

	sum := 0

	for _, h := range occupied {
		sum += h
	}

	return sum
}

func (v *Layout) getWidthAvailable() int {
	// occupied := v.getFooterHight() + v.getBodyHeight()
	return max(0, v.getWidthtWithinBorder())
}

func (v *Layout) getHeightWithinBorder() int {
	return max(0, v.size.GetHeight()-offsetBecauseBorder)
}

func (v *Layout) getWidthtWithinBorder() int {
	return max(0, v.size.GetWidth()-offsetBecauseBorder)
}

// Methods
func (v *Layout) makeVstack() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		v.makeHeaderView(),
		v.makeSectionView(),
		v.makeFooterView(),
	)
}

func (v *Layout) makeHstack() string {
	return v.makeSidepanelView()
}

func (v *Layout) makeView() string {
	vstack := v.makeVstack()
	hstack := v.makeHstack()

	mainWidth := max(0, v.getWidthAvailable()-lipgloss.Width(hstack))

	vstack = lipgloss.NewStyle().Width(mainWidth).Render(vstack)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		vstack,
		hstack,
	)
}

// Helpers
func trimParagraphPadding(content string) string {
	return strings.TrimRight(content, "\n")
}
