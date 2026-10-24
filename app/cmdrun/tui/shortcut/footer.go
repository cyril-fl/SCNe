package shortcut

// i19n.Localize("cmd.ctrl+n:new"

var (
	EDIT = NewShortcut(Template{
		Command: CtrlEnter,
		Label:   "edit",
	})
	NEW = NewShortcut(Template{
		Command: CtrlN,
		Label:   "new",
	})
	OPEN = NewShortcut(Template{
		Command: CtrlEnter,
		Label:   "open",
	})
	QUIT = NewShortcut(Template{
		Command: CtrlQ,
		Label:   "quit",
	})
	SAVEFORM = NewShortcut(Template{
		Command: CtrlEnter,
		Label:   "save:form",
	})
	SAVE = NewShortcut(Template{
		Command: CtrlS,
		Label:   "save",
	})
	SELECT = NewShortcut(Template{
		Command: Enter,
		Label:   "select",
	})
	SEND = NewShortcut(Template{
		Command: Enter,
		Label:   "send",
	})
	THROW = NewShortcut(Template{
		Command: CtrlEnter,
		Label:   "roll",
	})
)
