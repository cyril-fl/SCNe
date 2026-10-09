package layout

import (
	"solopg/app/shared/components/layout/sidepanel"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (v *Layout) SetSidepanel(parts ...string) {
	v.sidepanel = sidepanel.New(parts...)
}

func (v *Layout) getSidepanelHeight() int {
	if v.sidepanel != nil {
		return v.sidepanel.GetHeight()
	}
	return 0
}
func (v *Layout) MergeSidepanel(other tea.Model) {
	if v.sidepanel != nil {
		v.sidepanel.Merge(other)
	}
}

func (v *Layout) makeSidepanelView(view *strings.Builder) {
	if v.sidepanel != nil {
		view.WriteString(v.spaceSidepanel())
		view.WriteString(trimParagraphPadding(v.sidepanel.View()))
	}
}

func (v *Layout) spaceSidepanel() string {
	maxHeight := v.getHeightAvailable()
	return strings.Repeat("\n", max(0, maxHeight))
}
