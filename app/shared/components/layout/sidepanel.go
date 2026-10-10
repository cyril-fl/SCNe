package layout

import (
	"solopg/app/shared/components/layout/sidepanel"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (v *Layout) SetSidepanel(parts ...string) {
	if v.sidepanel == nil {
		v.sidepanel = sidepanel.Nil()
	}

	v.sidepanel = sidepanel.New(parts...)
}

func (v *Layout) getSidepanelHeight() int {
	if v.sidepanel == nil {
		return 0
	}
	
	return v.sidepanel.GetHeight()
}

func (v *Layout) MergeSidepanel(other tea.Model) {
	if v.sidepanel == nil {
		v.sidepanel = sidepanel.Nil()
	}

	v.sidepanel.Merge(other)
}


func (v *Layout) makeSidepanelView() string {
	if v.sidepanel == nil {
		return ""
	}

	view := strings.Builder{}
	view.WriteString(v.spaceSidepanel())
	view.WriteString(trimParagraphPadding(v.sidepanel.View()))

	return view.String()
}

func (v *Layout) spaceSidepanel() string {
	if v.sidepanel == nil {
		return ""
	}

	remainingHeight := v.getHeightAvailable() - v.getSidepanelHeight()
	return strings.Repeat("\n", max(0, remainingHeight))
}
