package direction

import (
	sharedtui "solopg/app/shared/tui"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
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

	switch sharedtui.EventKey(key.String()) {
	case sharedtui.KEY_UP:
		if isFirstEl := currentIndex == 0; isFirstEl && previousIndex == 0 {
			direction = PREVIOUS
		}
	case sharedtui.KEY_DOWN:
		if isLastEL := currentIndex == itemCount-1; isLastEL && previousIndex == currentIndex {
			direction = NEXT
		}
	}

	return updatedList, direction
}
