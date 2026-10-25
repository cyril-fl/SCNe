package layout

import (
	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/layout/sidepanel"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (v *Layout) SetSidepanel(parts ...string) {
	if v.sidepanel == nil {
		v.sidepanel = sidepanel.Nil()
	}

	v.sidepanel = sidepanel.New(parts...)
}

func (v *Layout) getSidepanelWidth() int {
	if v.sidepanel == nil {
		return 0
	}
	return v.sidepanel.GetWidth()
}

func (v *Layout) getSidepanelHeight() int {
	if v.sidepanel == nil {
		return 0
	}

	return v.sidepanel.GetHeight()
}

func (v *Layout) getAvailablePanelSize() size.Size {
	panelsize := 34
	return size.New(size.Template{
		Width:  panelsize,
		Height: v.getHeightWithinBorder(),
	})
}

func (v *Layout) MergeSidepanel(other tea.Model) {
	if v.sidepanel == nil {
		v.sidepanel = sidepanel.Nil()
	}

	v.sidepanel.SetSize(other, v.getAvailablePanelSize())
	v.sidepanel.Merge(other)
}

func (v *Layout) makeSidepanelView() string {
	if v.sidepanel == nil {
		return ""
	}

	view := strings.Builder{}
	view.WriteString(trimParagraphPadding(v.sidepanel.View()))

	return view.String()
}
