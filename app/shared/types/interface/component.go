package interfass

import (
	"solopg/app/cmdrun/types/size"

	tea "charm.land/bubbletea/v2"
)

type Amandable interface {
	Amend(content string)
}

type Viewable interface {
	GetView() string
	Update(msg tea.Msg) tea.Cmd
	SetSize(size size.Size)
}
