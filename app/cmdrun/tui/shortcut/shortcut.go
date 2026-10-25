package shortcut

import (
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/style"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

type Command string

const (
	CtrlQ      Command = "^Q"
	CtrlS      Command = "^S"
	CtrlN      Command = "^N"
	CtrlE      Command = "^E"
	ShiftEnter Command = "⇧↳"
	CtrlEnter  Command = "^↳"
	Enter      Command = "↳"
	AltArrows  Command = "⌥↑↓"
	CtrlPlus   Command = "^+"
	Esc        Command = "⎋"
	Alt        Command = "⌥"
)

type Template struct {
	Command Command
	Label   string
}

func NewShortcut(params Template) Shortcut {
	return Shortcut{
		command: params.Command,
		label:   params.Label,
	}
}

type Shortcut struct {
	command Command
	label   string
}

func (a Shortcut) String() string {
	return style.Key.Render(string(a.command)) + " " + style.Muted.Render(plumtransform.Capitalize(i19n.Localize(a.label)))
}
