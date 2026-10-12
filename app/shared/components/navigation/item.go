package navigation

import (
	"solopg/app/cmdrun/types/direction"
	"solopg/app/cmdrun/types/size"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// - Contract - //
// TODO verrifier ce qui est vraiment utilise et commain
// Remplacer le reste pas d'autre interface
type Item interface {
	ID() string

	SetFocus(bool)

	GetList() *list.Model
	SetList(*list.Model)

	SetSize(size size.Size)

	GetMenuView() string
	GetView() string
	GetFooter() []string

	HandleUpdate(tea.Msg) tea.Cmd
}

// - Context - //
type Context struct {
	Menu             Item
	CurrentMenuIndex int
	SibblingCount    int
}

func (m *Context) IsLastMenuElement() bool {
	return m.CurrentMenuIndex == m.SibblingCount-1
}

func (m *Context) IsFirstMenuElement() bool {
	return m.CurrentMenuIndex == 0
}

// - Update - //

// - Helper - //
func HandleKeyArrow(metadata Context, key tea.KeyPressMsg) direction.Direction {
	list, newdirection := direction.GetListDirection(metadata.Menu.GetList(), key)

	metadata.Menu.SetList(&list)

	isFirstElementToNext := metadata.IsFirstMenuElement() && newdirection == direction.NEXT
	isLastElementToPrevious := metadata.IsLastMenuElement() && newdirection == direction.PREVIOUS
	isBetweenElements := !(metadata.IsFirstMenuElement() || metadata.IsLastMenuElement())

	if isFirstElementToNext || isBetweenElements || isLastElementToPrevious {
		if menu, ok := metadata.Menu.(Openable); ok {
			menu.SetOpen(false)
		}
	}

	return newdirection
}
