package formpage

import (
	"solopg/app/cmdrun/tui/shortcut"
	"solopg/app/cmdrun/types/id"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/models"
	"solopg/app/shared/types/primitive"

	tea "charm.land/bubbletea/v2"
)

// Item is an entry displayed by a FormPage. Values are used to initialise the
// form when the entry is opened for editing. ID is passed back to Save.
type Item struct {
	ID          id.ID
	Name        string
	Description string
	Values      map[string]string
	Locked      bool
}

// Template configures a Page.
type Template struct {
	Items   func() []Item
	NewForm func(values map[string]string) *form.Form
	// Save receives a zero ID for a new entry.
	Save func(entryID id.ID, values map[string]string) error
	// Action is optional.
	Action *Action
}

// Action runs on the selected entry with Ctrl+Enter while the list is shown.
type Action struct {
	Label string
	Run   func(entryID id.ID)
}

// Page shows a list of items and opens a form to create or edit one.
type Page struct {
	primitive.List

	items   func() []Item
	newForm func(map[string]string) *form.Form
	save    func(id.ID, map[string]string) error
	action  *Action

	form    *form.Form
	entryID id.ID
	focused bool
}

// TODO regardader pour adopter cette pholosipheio ou c'est possible

// New returns a page listing template.Items.
func New(template Template) *Page {
	p := &Page{
		items:   template.Items,
		newForm: template.NewForm,
		save:    template.Save,
		action:  template.Action,
		focused: true,
	}

	p.intList()
	p.refresh()

	return p
}

// TODO la page a un décallage et ne s'affiche pas nickel
// Getters & Setters

// GetView renders the open form, or the list when none is open.
func (p *Page) GetView() string {
	if p.form != nil {
		return p.form.View().Content
	}
	li := p.GetList()
	return li.View()
}

// GetForm returns the open form, or nil.
func (p *Page) GetForm() *form.Form {
	return p.form
}

// GetFooter returns the shortcuts available in the current state.
func (p *Page) GetFooter() []string {
	isFormOpen := p.IsFormOpen()
	isLocked := p.isSelectedLocked()
	hasAction := p.action != nil

	switch {
	case isFormOpen && p.form.CanRepeat():
		return []string{shortcut.SAVEFORM.String(), shortcut.ADD.String()}
	case isFormOpen:
		return []string{shortcut.SAVEFORM.String()}
	case isLocked && hasAction:
		return []string{shortcut.NEW.String(), p.getActionFooter()}
	case isLocked:
		return []string{shortcut.NEW.String()}
	case hasAction:
		return []string{shortcut.NEW.String(), shortcut.EDIT.String(), p.getActionFooter()}
	default:
		return []string{shortcut.NEW.String(), shortcut.EDIT.String()}
	}
}

func (p *Page) getActionFooter() string {
	return shortcut.NewShortcut(shortcut.Template{
		Command: shortcut.CtrlEnter,
		Label:   p.action.Label,
	}).String()
}

func (p *Page) isSelectedLocked() bool {
	li := p.GetList()
	item := li.SelectedItem()

	if item == nil {
		return false
	}

	if item, ok := item.(models.Item[Item]); ok && !item.Value().Locked {
		return false
	}

	return true
}

// SetFocus sets the list focus.
func (p *Page) SetFocus(focused bool) {
	p.focused = focused
	models.SetListFocus(p.GetList(), focused)
}

// SetSize sets the list width; its height fits the content.
func (p *Page) SetSize(s size.Size) {
	list := p.GetList()
	p.SetListSize(size.Template{
		Width:  s.GetWidth(),
		Height: models.ListContentHeight(list),
	})
}

// IsFormOpen reports whether a form is open.
func (p *Page) IsFormOpen() bool {
	return p.form != nil
}

// Methods

// Update validates and saves the form on form.Validate, and routes other messages
// to the open form or the list.
func (p *Page) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case form.Validate:
		return p.handleValidate()
	case tea.KeyPressMsg:
		return p.handleKeyPress(msg)
	default:
		return p.handleModelUpdate(msg)
	}
}

// Close discards the open form.
func (p *Page) Close() {
	p.form = nil
	p.entryID = id.ID{}
}

// EditNewItem opens an empty form for a new entry.
func (p *Page) EditNewItem() {
	p.entryID = id.ID{}
	p.form = p.newForm(nil)
}

// EditItem opens the form for the selected entry.
// It returns false when nothing is selected or the entry is locked.
func (p *Page) EditItem() bool {
	list := p.GetList()

	selected, ok := list.SelectedItem().(models.Item[Item])
	if !ok {
		return false
	}

	entry := selected.Value()
	p.entryID = entry.ID
	if entry.Locked {
		return false
	}

	p.form = p.newForm(entry.Values)
	return true
}

// HandleEscape closes the open form and reports whether one was open.
func (p *Page) HandleEscape() bool {
	if p.form != nil {
		p.form = nil
		p.entryID = id.ID{}
		return true
	}
	return false
}
