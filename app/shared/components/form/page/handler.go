package formpage

import (
	"scne/app/cmdrun/types/id"
	"scne/app/shared/components/form"
	"scne/app/shared/components/models"
	"scne/app/shared/messages"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

func (p *Page) handleKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	if p.form != nil {
		return p.handleFormKeyPress(msg)
	}
	return p.handleListKeyPress(msg)
}

func (p *Page) handleFormKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	switch plumtypes.Key(msg.String()) {
	case plumtypes.KEY_ESC:
		p.HandleEscape()
		return nil
	case plumtypes.CMD_CTRL_ENTER:
		return form.Send[form.Validate]()
	default:
		return p.handleModelUpdate(msg)
	}
}

func (p *Page) handleListKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	switch plumtypes.Key(msg.String()) {
	case plumtypes.KEY_ESC:
		return p.handleKeyEsc()
	case plumtypes.SHORT_CTRL_N:
		return p.handleCtrlN()
	case "ctrl+e":
		return p.handleCtrlE()
	case plumtypes.CMD_CTRL_ENTER:
		return p.handleCtrlEnter()
	default:
		return p.handleModelUpdate(msg)
	}
}

func (p *Page) handleModelUpdate(msg tea.Msg) tea.Cmd {
	if p.form != nil {
		_, cmd := p.form.Update(msg)
		return cmd
	}

	list := p.GetList()
	newlist, cmd := list.Update(msg)
	p.SetList(&newlist)
	return cmd
}

func (p *Page) handleValidate() tea.Cmd {
	if p.form == nil {
		return nil
	}

	p.form.Submit()
	if p.form.HasErr() {
		return messages.Fail(p.form.GetErr())
	}
	if err := p.save(p.entryID, p.form.GetValuesAsMappedString()); err != nil {
		p.form.SetErr(err)
		return messages.Fail(err)
	}
	p.form = nil
	p.entryID = id.ID{}
	p.refresh()
	return nil
}

func (p *Page) handleKeyEsc() tea.Cmd {
	return form.Send[form.Close]()
}

func (p *Page) handleCtrlN() tea.Cmd {
	if p.form == nil {
		p.EditNewItem()
	}

	return nil
}

func (p *Page) handleCtrlEnter() tea.Cmd {
	if p.action == nil {
		return nil
	}

	list := p.GetList()
	selected, ok := list.SelectedItem().(models.Item[Item])
	if !ok {
		return nil
	}

	p.action.Run(selected.Value().ID)
	return nil
}

func (p *Page) handleCtrlE() tea.Cmd {
	if p.form == nil {
		p.EditItem()
	}
	return nil
}

func (p *Page) intList() {
	l := p.GetList()
	*l = models.NewList(nil, 0, 0)
}

func (p *Page) refresh() {
	l := p.GetList()
	items := p.items()
	listItems := make([]list.Item, 0, len(items))
	for _, item := range items {
		listItems = append(listItems, models.NewItem(item.Name, item.Description, item))
	}
	l.SetItems(listItems)
	models.SetListFocus(l, p.focused)
}
