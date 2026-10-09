package layout

import (
	"solopg/app/shared/components/layout/header"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (v *Layout) SetHeader(parts ...string) {
	v.header = header.New(parts...)
}

func (v *Layout) getHeaderHight() int {
	if v.header != nil {
		return v.header.GetHeight()
	}
	return 0
}
func (v *Layout) MergeHeader(other tea.Model) {
	if v.header != nil {
		v.header.Merge(other)
	}
}

func (v *Layout) makeHeaderView(view *strings.Builder) {
	if v.header != nil {
		view.WriteString(v.spaceHeader())
		view.WriteString(trimParagraphPadding(v.header.View()))
	}
}

func (v *Layout) spaceHeader() string {
	maxHeight := v.getHeightAvailable()
	return strings.Repeat("\n", max(0, maxHeight))
}
