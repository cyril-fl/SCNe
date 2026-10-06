package codexmenu

import (
	"solopg/app/cmdrun/domain/gameplay/codex"
	"solopg/app/cmdrun/tui/models"
	"solopg/app/cmdrun/tui/view/sidemenu"
	"solopg/app/shared/components/form"
	sharedtui "solopg/app/shared/tui"

	"solopg/app/cmdrun/types/direction"
	"solopg/app/cmdrun/types/size"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// - CodexMenu - //
// Menu
func NewSideMenu(params CodexMenuParams, focused bool) *codexMenu {
	items := makeCodexList(params)
	menu := list.New(items, list.NewDefaultDelegate(), params.Size.Width, params.Size.Height)
	models.ConfigureList(&menu)
	models.SetListFocus(&menu, focused)

	return &codexMenu{
		id:      "codex",
		focused: focused,
		list:    menu,
	}
}

type CodexMenuParams struct {
	Size  size.Size
	Codex *codex.Codex
}

type codexMenu struct {
	id      string
	focused bool
	list    list.Model
}

// / Getters and Setters
func (m *codexMenu) ID() string {
	return m.id
}

func (m *codexMenu) IsOpen() bool {
	for _, item := range m.list.Items() {
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

	for _, item := range m.list.Items() {
		item, ok := item.(models.Item[*codexMenuItem])
		if page := item.Value(); ok {
			page.setOpen(false)
		}
	}
}

func (m *codexMenu) SetFocus(focused bool) {
	m.focused = focused
	models.SetListFocus(&m.list, focused)
}

func (m *codexMenu) GetList() *list.Model {
	return &m.list
}

func (m *codexMenu) SetList(list list.Model) {
	m.list = list
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
func (m *codexMenu) HandleUpdate(params sidemenu.UpdateParams) (tea.Model, tea.Cmd) {
	switch msg := params.Msg.(type) {
	case form.NextField, form.PreviousField:
		return m.delegateInputToForm(params)
	case form.Validate:
		return m.handleFormSubmit(params)
	case form.Error:
		return m.handleFormError(params)
	case form.Scroll:
		return m.handleFormScroll(params)

	case tea.KeyPressMsg:
		switch msg.String() {
		case sharedtui.KEY_ESC:
			return m.handleKeyEsc(params)
		case sharedtui.KEY_UP, sharedtui.KEY_DOWN:
			return m.delegateInputToForm(params)
		case sharedtui.KEY_ENTER:
			return m.delegateInputToForm(params)
		case sharedtui.CMD_SHIFT_ENTER, sharedtui.CMD_ALT_ENTER:
			return m.handleKeyShiftEnter(params)
		default:
			return m.delegateInputToForm(params)
		}
	}

	return params.Delegate(params.Msg)
}

func (m *codexMenu) HandleDirectionInput(direction direction.Direction) {
	_ = direction
	m.updatePageOnRedirection()
}

/*
NOT-IMPLEMENTED-YET
func (m *codexMenu) HandleScrollInput(viewport viewport.Model) viewport.Model {
 	return m.handleFormScrollInput(viewport)
}
*/

// Items
type codexMenuItem struct {
	id     string
	isOpen bool
	table  codex.Table

	form     *form.Form
	newForm  func() *form.Form
	showForm bool
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
