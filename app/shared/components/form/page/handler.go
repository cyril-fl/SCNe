package formpage

import (
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/models"
	"solopg/app/shared/messages"
	sharedtui "solopg/app/shared/tui"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

func (p *Page) handleKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	if p.form != nil {
		return p.handleFormKeyPress(msg)
	}
	return p.handleListKeyPress(msg)
}

func (p *Page) handleFormKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	switch sharedtui.EventKey(msg.String()) {
	case sharedtui.KEY_ESC:
		p.HandleEscape()
		return nil
	case sharedtui.CMD_ALT_ENTER:
		return form.Send[form.Validate]()
	default:
		return p.handleModelUpdate(msg)
	}
}

func (p *Page) handleListKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	switch sharedtui.EventKey(msg.String()) {
	case sharedtui.KEY_ESC:
		return p.handleKeyEsc()
	case sharedtui.SHORT_CTRL_N:
		return p.handleCtrlN()
	case sharedtui.CMD_ALT_ENTER:
		return p.handleAltEnter()
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
	if err := p.save(p.index, p.form.GetValuesAsMappedString()); err != nil {
		p.form.SetErr(err)
		return messages.Fail(err)
	}
	p.form = nil
	p.index = -1
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

func (p *Page) handleAltEnter() tea.Cmd {
	if p.form == nil {
		p.EditItem()
		return nil
	}
	return form.Send[form.Validate]()
}

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
