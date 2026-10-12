package layout

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/layout/section"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (v *Layout) SetSection(parts ...string) {
	if v.section == nil {
		v.section = section.Nil()
	}

	v.section = section.New(parts...)
}

func (v *Layout) getSectionHeight() int {
	if v.section == nil {
		return 0
	}

	return v.section.GetHeight()
}

func (v *Layout) MergeSection(other tea.Model) {
	if v.section == nil {
		v.section = section.Nil()
	}

	v.section.Merge(other)
	v.section.SetSize(other, v.getAvailableSectionSize())
}

func (v *Layout) getAvailableSectionSize() size.Size {
	return size.New(size.Template{
		Width:  v.size.GetWidth() - v.getSidepanelWidth(),
		Height: v.size.GetHeight() - v.getHeaderHeight() - v.getFooterHeight(),
	})
}

func (v *Layout) makeSectionView() string {
	if v.section == nil {
		return ""
	}

	view := strings.Builder{}
	// view.WriteString(v.spaceSection())
	view.WriteString(trimParagraphPadding(v.section.View()))

	return view.String()
}

func (v *Layout) spaceSection() string {
	if v.section == nil {
		return ""
	}
	maxHeight := v.getHeightAvailable()
	return strings.Repeat("\n", max(0, maxHeight))
}
