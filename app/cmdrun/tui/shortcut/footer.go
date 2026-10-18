package shortcut

// i19n.Localize("cmd.ctrl+n:new"

var (
	EDIT = NewShortcut(Template{
		Command: AltEnter,
		Label:   "edit",
	})
	NEW = NewShortcut(Template{
		Command: CtrlN,
		Label:   "new",
	})
	OPEN = NewShortcut(Template{
		Command: AltEnter,
		Label:   "open",
	})
	QUIT = NewShortcut(Template{
		Command: CtrlQ,
		Label:   "quit",
	})
	SAVEFORM = NewShortcut(Template{
		Command: AltEnter,
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
		Command: AltEnter,
		Label:   "roll",
	})
)
