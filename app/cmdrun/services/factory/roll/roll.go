package roll

import (
	sharedtui "solopg/app/shared/tui"

	tea "charm.land/bubbletea/v2"
)

type Roller interface {
	Roll() (string, error)
}

type RollResultMsg struct {
	Content string
}

func Execute(roller Roller) tea.Cmd {
	return func() tea.Msg {
		content, err := roller.Roll()
		if err != nil {
			return sharedtui.ErrorMsg{Err: err}
		}
		return RollResultMsg{Content: content}
	}
}
