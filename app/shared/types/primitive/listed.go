package primitive

import (
	"solopg/app/cmdrun/types/size"

	"charm.land/bubbles/v2/list"
)

// TEST
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

func (m *List) SetList(listModel *list.Model) {
	m.list = *listModel
}


func (m *List) SetListSize(template size.Template) {
		s := size.New(template)
		m.list.SetSize(s.GetWidth(), s.GetHeight())
}
