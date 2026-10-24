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

type Template struct {
	Items   func() []Item
	NewForm func(values map[string]string) *form.Form
	// Save receives a zero ID for a new entry.
	Save func(entryID id.ID, values map[string]string) error
}

type Page struct {
	primitive.List

	items   func() []Item
	newForm func(map[string]string) *form.Form
	save    func(id.ID, map[string]string) error

	form    *form.Form
	entryID id.ID
}

// TODO regardader pour adopter cette pholosipheio ou c'est possible
func New(template Template) *Page {
	p := &Page{
		items:   template.Items,
		newForm: template.NewForm,
		save:    template.Save,
	}

	p.intList()
	p.refresh()

	return p
}

// TODO la page a un décallage et ne s'affiche pas nickel
// Getters & Setters
func (p *Page) GetView() string {
	if p.form != nil {
		return p.form.View().Content
	}
	li := p.GetList()
	return li.View()
}

func (p *Page) GetForm() *form.Form {
	return p.form
}

func (p *Page) GetFooter() []string {
	isFormOpen := p.IsFormOpen()

	switch {
	case isFormOpen && p.form.CanRepeat():
		return []string{shortcut.SAVEFORM.String(), shortcut.ADD.String()}
	case isFormOpen:
		return []string{shortcut.SAVEFORM.String()}
	case p.isSelectedLocked():
		return []string{shortcut.NEW.String()}
	default:
		return []string{shortcut.NEW.String(), shortcut.EDIT.String()}
	}
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

func (p *Page) SetFocus(focused bool) {
	models.SetListFocus(p.GetList(), focused)
}

func (p *Page) SetSize(s size.Size) {
	list := p.GetList()
	p.SetListSize(size.Template{
		Width:  s.GetWidth(),
		Height: models.ListContentHeight(list),
	})
}

func (p *Page) IsFormOpen() bool {
	return p.form != nil
}

// Methods
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

func (p *Page) Close() {
	p.form = nil
	p.entryID = id.ID{}
}

func (p *Page) EditNewItem() {
	p.entryID = id.ID{}
	p.form = p.newForm(nil)
}

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

func (p *Page) HandleEscape() bool {
	if p.form != nil {
		p.form = nil
		p.entryID = id.ID{}
		return true
	}
	return false
}
