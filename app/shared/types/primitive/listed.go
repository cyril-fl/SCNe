package primitive

import "charm.land/bubbles/v2/list"

type Listed struct {
	list    list.Model
}

func NewListed(listModel list.Model) Listed {
	return Listed{
		list: listModel,
	}
}

func (m *Listed) GetList() *list.Model {
	return &m.list
}

func (m *Listed) SetList(list *list.Model) {
	m.list = *list
}