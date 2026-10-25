package primitive

import (
	"scne/app/cmdrun/types/size"

	"charm.land/bubbles/v2/list"
)

// List embeds a bubbles list into a model.
type List struct {
	list list.Model
}

// Listed wraps listModel.
func Listed(listModel list.Model) List {
	return List{
		list: listModel,
	}
}

// GetList returns the wrapped list, mutable in place.
func (m *List) GetList() *list.Model {
	return &m.list
}

// SetList replaces the wrapped list with a copy of listModel.
func (m *List) SetList(listModel *list.Model) {
	m.list = *listModel
}

// SetListSize resizes the list.
func (m *List) SetListSize(template size.Template) {
	s := size.New(template)
	m.list.SetSize(s.GetWidth(), s.GetHeight())
}
