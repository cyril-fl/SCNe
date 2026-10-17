package formpage

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/models"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/types/primitive"

	tea "charm.land/bubbletea/v2"
)

// Item is an entry displayed by a FormPage. Values are used to initialise the
// form when the entry is opened for editing.
type Item struct {
	Name        string
	Description string
	Values      map[string]string
	Locked      bool
}

type Template struct {
	Items   func() []Item
	NewForm func(values map[string]string) *form.Form
	Save    func(index int, values map[string]string) error
}

type Page struct {
	primitive.List

	items   func() []Item
	newForm func(map[string]string) *form.Form
	save    func(int, map[string]string) error

	form  *form.Form
	index int
}

// TODO regardader pour adopter cette pholosipheio ou c'est possible
func New(template Template) *Page {
	p := &Page{
		items:   template.Items,
		newForm: template.NewForm,
		save:    template.Save,
		index:   -1,
	}

	p.intList()
	p.refresh()

	return p
}

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
	if p.form != nil {
		return []string{i19n.Localize("cmd.alt+enter:validate")}
	}
	return []string{i19n.Localize("cmd.ctrl+n:new")}
}

func (p *Page) SetFocus(focused bool) {
	models.SetListFocus(p.GetList(), focused)
}

func (p *Page) SetSize(s size.Size) {
	list := p.GetList()
	p.SetListSize(size.Template{
		Width:  s.GetWidth(),
		Height:	models.ListContentHeight(list),
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
	p.index = -1
}

func (p *Page) EditNewItem() {
	p.index = -1
	p.form = p.newForm(nil)
}

func (p *Page) EditItem() bool {
	list := p.GetList()

	selected, ok := list.SelectedItem().(models.Item[Item])
	if !ok {
		return false
	}

	p.index = list.Index()
	entry := selected.Value()
	if entry.Locked {
		return false
	}

	p.form = p.newForm(entry.Values)
	return true
}

func (p *Page) HandleEscape() bool {
	if p.form != nil {
		p.form = nil
		p.index = -1
		return true
	}
	return false
}
