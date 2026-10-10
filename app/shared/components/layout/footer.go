package layout

import (
	"solopg/app/shared/components/layout/footer"
	"strings"

	tea "charm.land/bubbletea/v2"
)



func (v *Layout) SetFooter(parts ...string) {
	if v.footer == nil {
		v.footer = footer.Nil()
	}

	v.footer = footer.New(parts...)
}

func (v *Layout) getFooterHeight() int {
	if v.footer != nil {
		return v.footer.GetHeight()
	}
	
	return 0
}

func (v *Layout) MergeFooter(other tea.Model) {
	if v.footer == nil {
		v.footer = footer.Nil()
	}
	
	v.footer.Merge(other)
}

func (v *Layout) makeFooterView() string {
	if v.footer == nil {
		return ""
	}
	
	view := strings.Builder{}
	view.WriteString(v.spaceFooter())
	view.WriteString(trimParagraphPadding(v.footer.View()))
		
	return view.String()
}

func (v *Layout) spaceFooter() string {
	if v.footer == nil {
		return ""
	}
	maxHeight := v.getHeightAvailable()
	return strings.Repeat("\n", max(0, maxHeight))
}