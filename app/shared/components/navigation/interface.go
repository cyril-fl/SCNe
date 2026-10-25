package navigation

import (
	"solopg/app/cmdrun/types/size"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// SectionInterface is a navigation block.
type SectionInterface interface {
	ID() string

	SetSize(size size.Size)

	GetView() string
	GetFooter() []string

	HandleUpdate(tea.Msg) tea.Cmd
}

// FocusableSectionInterface is a navigation block with a selectable list.
type FocusableSectionInterface interface {
	SectionInterface

	SetFocus(bool)

	GetList() *list.Model
	SetList(*list.Model)
}
