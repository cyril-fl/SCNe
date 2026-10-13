package primitive

import "charm.land/bubbles/v2/list"

type List struct {
	list list.Model
}

func Listed(listModel list.Model) List {
	return List{
		list: listModel,
	}
}

func (m *List) GetList() *list.Model {
	return &m.list
}

func (m *List) SetList(list *list.Model) {
	m.list = *list
}
