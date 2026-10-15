package formpage

import (
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/models"
	"solopg/app/shared/messages"
	"solopg/app/shared/services/i19n"
	sharedtui "solopg/app/shared/tui"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// Item is an entry displayed by a FormPage. Values are used to initialise the
// form when the entry is opened for editing.
type Item struct {
	Title       string
	Description string
	Values      map[string]string
	Locked      bool
}

type Template struct {
	Items   func() []Item
	NewForm func(values map[string]string) *form.Form
	Save    func(index int, values map[string]string) error
}

type CloseMsg struct{}

type FormPage struct {
	items   func() []Item
	newForm func(map[string]string) *form.Form
	save    func(int, map[string]string) error
	list    list.Model
	form    *form.Form
	index   int
}

// Page is kept as a concise alias for callers that prefer the component's
// shorter name.
type Page = FormPage

func New(template Template) *FormPage {
	p := &FormPage{
		items:   template.Items,
		newForm: template.NewForm,
		save:    template.Save,
		index:   -1,
	}
	p.list = list.New(nil, list.NewDefaultDelegate(), 0, 0)
	models.ConfigureList(&p.list)
	p.Refresh()
	return p
}

func (p *FormPage) Refresh() {
	items := p.items()
	listItems := make([]list.Item, 0, len(items))
	for _, item := range items {
		listItems = append(listItems, models.NewItem(item.Title, item.Description, item))
	}
	p.list.SetItems(listItems)
}

func (p *FormPage) SetSize(s size.Size) {
	p.list.SetSize(s.GetWidth(), s.GetHeight())
}

func (p *FormPage) SetFocus(focused bool) {
	models.SetListFocus(&p.list, focused)
}

func (p *FormPage) GetList() *list.Model { return &p.list }

func (p *FormPage) IsFormOpen() bool { return p.form != nil }

func (p *FormPage) Form() *form.Form { return p.form }

func (p *FormPage) OpenNew() {
	p.index = -1
	p.form = p.newForm(nil)
}

func (p *FormPage) OpenSelected() bool {
	selected, ok := p.list.SelectedItem().(models.Item[Item])
	if !ok {
		return false
	}
	p.index = p.list.Index()
	entry := selected.Value()
	if entry.Locked {
		return false
	}
	p.form = p.newForm(entry.Values)
	return true
}

// HandleEscape closes the form first. If no form is open, it reports that the
// owning menu may close the page.
func (p *FormPage) HandleEscape() bool {
	if p.form != nil {
		p.form = nil
		p.index = -1
		return true
	}
	return false
}

func (p *FormPage) Close() {
	p.form = nil
	p.index = -1
}

func (p *FormPage) Update(msg tea.Msg) tea.Cmd {
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
			p.Refresh()
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
			p.OpenNew()
			return nil
		case sharedtui.CMD_ALT_ENTER:
			p.OpenSelected()
			return nil
		}
	}
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return cmd
}

func (p *FormPage) View() string {
	if p.form != nil {
		return p.form.View().Content
	}
	return p.list.View()
}

func (p *FormPage) GetView() string { return p.View() }

func (p *FormPage) GetFooter() []string {
	if p.form != nil {
		return []string{i19n.Localize("cmd.alt+enter:validate")}
	}
	return []string{i19n.Localize("cmd.ctrl+n:new")}
}
