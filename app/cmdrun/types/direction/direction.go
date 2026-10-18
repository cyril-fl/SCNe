package direction

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

type Direction string

const (
	PREVIOUS Direction = "previous"
	NEXT     Direction = "next"
	NONE     Direction = "none"
)

func GetListDirection(list *list.Model, key tea.KeyPressMsg) (list.Model, Direction) {
	direction := NONE
	itemCount := len(list.Items())

	previousIndex := list.Index()
	updatedList, _ := list.Update(key)
	currentIndex := list.Index()

	switch plumtypes.Key(key.String()) {
	case plumtypes.KEY_UP:
		if isFirstEl := currentIndex == 0; isFirstEl && previousIndex == 0 {
			direction = PREVIOUS
		}
	case plumtypes.KEY_DOWN:
		if isLastEL := currentIndex == itemCount-1; isLastEL && previousIndex == currentIndex {
			direction = NEXT
		}
	}

	return updatedList, direction
}
