package layout

import (
	"scne/app/shared/components/layout/header"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (v *Layout) getHeaderHeight() int {
	if v.header == nil {
		return 0
	}

	return v.header.GetHeight()
}

func (v *Layout) SetHeader(other tea.Model) {
	if v.header == nil {
		v.header = header.Nil()
	}

	v.header.SetHeader(other)
}

func (v *Layout) makeHeaderView() string {
	if v.header == nil {
		return ""
	}

	view := strings.Builder{}
	view.WriteString(trimParagraphPadding(v.header.View()))
	// view.WriteString(v.spaceHeader())

	return view.String()
}

func (v *Layout) spaceHeader() string {
	if v.header == nil {
		return ""
	}

	maxHeight := v.getHeightAvailable()
	return strings.Repeat("\n", max(0, maxHeight))
}
