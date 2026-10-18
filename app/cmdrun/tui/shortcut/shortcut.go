package shortcut

import (
	"fmt"
	"solopg/app/shared/services/i19n"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

type Command string

const (
	CtrlQ      Command = "^Q"
	CtrlS      Command = "^S"
	CtrlN      Command = "^N"
	ShiftEnter Command = "⇧↳"
	AltEnter   Command = "⌥↳"
	Enter      Command = "↳"
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
	return fmt.Sprintf("[%s] %s", a.command, plumtransform.Capitalize(i19n.Localize(a.label)))
}
