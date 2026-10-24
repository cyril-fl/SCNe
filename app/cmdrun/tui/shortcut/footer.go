package shortcut

// i19n.Localize("cmd.ctrl+n:new"

var (
	ADD = NewShortcut(Template{
		Command: CtrlPlus,
		Label:   "add",
	})
	EDIT = NewShortcut(Template{
		Command: CtrlE,
		Label:   "edit",
	})
	NAVIGATE = NewShortcut(Template{
		Command: AltArrows,
		Label:   "navigate",
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
