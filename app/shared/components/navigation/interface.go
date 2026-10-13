package navigation

import (
	"solopg/app/cmdrun/types/size"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type SectionInterface interface {
	ID() string

	SetFocus(bool)

	GetList() *list.Model
	SetList(*list.Model)

	SetSize(size size.Size)

	GetView() string
	GetFooter() []string

	HandleUpdate(tea.Msg) tea.Cmd
}
