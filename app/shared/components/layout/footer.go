package layout

import (
	"solopg/app/shared/components/layout/footer"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (v *Layout) SetFooter(parts ...string) {
	v.footer = footer.New(parts...)
}

func (v *Layout) getFooterHight() int {
	if v.footer != nil {
		return v.footer.GetHeight()
	}
	return 0
}
func (v *Layout) MergeFooter(other tea.Model) {
	if v.footer != nil {
		v.footer.Merge(other)
	}
}

func (v *Layout) makeFooterView(view *strings.Builder) {
	if v.footer != nil {
		view.WriteString(v.spaceFooter())
		view.WriteString(trimParagraphPadding(v.footer.View()))
	}
}

func (v *Layout) spaceFooter() string {
	maxHeight := v.getHeightAvailable()
	return strings.Repeat("\n", max(0, maxHeight))
}
