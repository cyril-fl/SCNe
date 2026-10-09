package layout

import (
	"solopg/app/shared/components/layout/section"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (v *Layout) SetSection(parts ...string) {
 v.section = section.New(parts...)
}

func (v *Layout) getSectionHeight() int {
	if v.section != nil {
		return v.section.GetHeight()
	}
	return 0
}

func (v *Layout) MergeSection(other tea.Model) {
	if v.section != nil {
		v.section.Merge(other)
	}
}

func (v *Layout) makeSectionView(view *strings.Builder) {
	if v.section != nil {
		view.WriteString(v.spaceSection())
		view.WriteString(trimParagraphPadding(v.section.View()))
	}
}

func (v *Layout) spaceSection() string {
	maxHeight := v.getHeightAvailable()
	return strings.Repeat("\n", max(0, maxHeight))
}