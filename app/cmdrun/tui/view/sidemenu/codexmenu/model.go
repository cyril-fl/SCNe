package codexmenu

import (
	"solopg/app/cmdrun/domain/gameplay/codex"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/models"
	"solopg/app/shared/components/navigation"
	sharedtui "solopg/app/shared/tui"

	"solopg/app/cmdrun/types/direction"
	"solopg/app/cmdrun/types/size"

	tea "charm.land/bubbletea/v2"
)

// - CodexMenu - //
// Menu
func NewSideMenu(params CodexMenuParams, focused bool) *codexMenu {
	return &codexMenu{
		Menu: navigation.NewMenu(navigation.TemplateMenu[*codexMenuItem]{
			ID:      "codex",
			Options: makeCodexList(params),
			Focused: focused,
		}),
	}
}

type CodexMenuParams struct {
	Size  size.Size
	Codex *codex.Codex
}

type codexMenu struct {
	*navigation.Menu
}

// / Getters and Setters
func (m *codexMenu) IsOpen() bool {
	for _, item := range m.GetList().Items() {
		if page, ok := item.(models.Item[*codexMenuItem]); !ok {
			continue
		} else if page.Value().isOpen {
			return true
		}
	}
	return false
}

func (m *codexMenu) SetOpen(open bool) {
	if open {
		return
	}

	for _, item := range m.GetList().Items() {
		item, ok := item.(models.Item[*codexMenuItem])
		if page := item.Value(); ok {
			page.setOpen(false)
		}
	}
}

func (menu *codexMenu) GetFormFromCurrentPage() *form.Form {
	if currentPage := menu.getCurrentPage(); currentPage != nil && menu.IsOpen() {
		return currentPage.form
	}
	return nil
}

func (menu *codexMenu) SetFormOnCurrentPage(form form.Form) {
	if currentPage := menu.getCurrentPage(); currentPage != nil && menu.IsOpen() {
		currentPage.form = &form
	}
}

// / Handlers
func (m *codexMenu) HandleUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case form.NextField, form.PreviousField:
		return m.delegateInputToForm(msg)
	case form.Validate:
		return m.handleFormSubmit(msg)
	case form.Error:
		return m.handleFormError(msg)
	// case form.Scroll:
	// 	return m.handleFormScroll(msg)

	case tea.KeyPressMsg:
		switch sharedtui.EventKey(msg.String()) {
		case sharedtui.KEY_ESC:
			return m.handleKeyEsc(msg)
		case sharedtui.KEY_UP, sharedtui.KEY_DOWN:
			return m.delegateInputToForm(msg)
		case sharedtui.KEY_ENTER:
			return m.delegateInputToForm(msg)
		case sharedtui.CMD_SHIFT_ENTER, sharedtui.CMD_ALT_ENTER:
			return m.handleKeyShiftEnter(msg)
		default:
			return m.delegateInputToForm(msg)
		}
	}

	return nil
}

func (m *codexMenu) HandleDirectionInput(direction direction.Direction) {
	_ = direction
	m.updatePageOnRedirection()
}

// Items
type codexMenuItem struct {
	id     string
	isOpen bool
	table  codex.Table

	form     *form.Form
	newForm  func() *form.Form
	showForm bool
}

func (m *codexMenuItem) GetName() string {
	return m.id
}

type sideMenuItemsParams struct {
	id    string
	table codex.Table
	form  func() *form.Form
}

func newMenuItem(params sideMenuItemsParams) *codexMenuItem {
	return &codexMenuItem{
		id:       params.id,
		isOpen:   false,
		table:    params.table,
		newForm:  params.form,
		showForm: false,
	}
}

// Message
func sendMsg() tea.Cmd {
	return func() tea.Msg {
		return Msg{}
	}
}

type Msg struct {
}
