package formpage

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/models"
	"solopg/app/shared/messages"
	"solopg/app/shared/services/i19n"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/types/primitive"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// Item is an entry displayed by a FormPage. Values are used to initialise the
// form when the entry is opened for editing.
type Item struct {
	Name       string
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

	form    *form.Form
	index   int
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
	p.SetSize(s)
}

func (p *Page) IsFormOpen() bool { 
	return p.form != nil 
}


// Methods
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


func (p *Page) Update(msg tea.Msg) tea.Cmd {
	list := p.GetList()

	if p.form != nil {
		switch value := msg.(type) {
		case form.Validate:
			p.form.Submit()
			if p.form.HasErr() {
				return messages.Fail(p.form.GetErr())
			}
			if err := p.save(p.index, p.form.GetValuesAsMappedString()); err != nil {
				p.form.SetErr(err)
				return messages.Fail(err)
			}
			p.form = nil
			p.index = -1
			p.refresh()
			return nil
		case tea.KeyPressMsg:
			switch sharedtui.EventKey(value.String()) {
			case sharedtui.KEY_ESC:
				p.HandleEscape()
				return nil
			case sharedtui.CMD_ALT_ENTER:
				return form.Send[form.Validate]()
			}
		}
		_, cmd := p.form.Update(msg)
		return cmd
	}

	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch sharedtui.EventKey(key.String()) {
		case sharedtui.KEY_ESC:
			return func() tea.Msg { return CloseMsg{} }
		case sharedtui.SHORT_CTRL_N:
			p.EditNewItem()
			return nil
		case sharedtui.CMD_ALT_ENTER:
			p.EditItem()
			return nil
		}
	}

	newlist, cmd := list.Update(msg)
	p.SetList(newlist)
	return cmd
}

// Methods
func (p *Page) intList() {
	l := p.GetList()
	*l = list.New(nil, list.NewDefaultDelegate(), 0, 0)
	models.ConfigureList(l)
}

func (p *Page) refresh() {
	l := p.GetList()
	items := p.items()
	listItems := make([]list.Item, 0, len(items))
	for _, item := range items {
		listItems = append(listItems, models.NewItem(item.Name, item.Description, item))
	}
	l.SetItems(listItems)
}

// Helpers
type CloseMsg struct{}

// Implementations

// HandleEscape closes the form first. If no form is open, it reports that the
// owning menu may close the page.
func (p *Page) HandleEscape() bool {
	if p.form != nil {
		p.form = nil
		p.index = -1
		return true
	}
	return false
}