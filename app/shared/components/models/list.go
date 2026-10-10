package models

import (
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
)

type Item[T any] struct {
	title       string
	description string
	value       T
}

func NewItem[T any](title, description string, value T) Item[T] {
	return Item[T]{
		title:       title,
		description: description,
		value:       value,
	}
}

func (i Item[T]) Title() string       { return i.title }
func (i Item[T]) Description() string { return i.description }
func (i Item[T]) FilterValue() string { return i.title }
func (i Item[T]) Value() T            { return i.value }

func ConfigureList(menu *list.Model) {
	menu.Title = ""
	menu.SetShowFilter(false)
	menu.SetShowPagination(false)
	menu.SetShowHelp(false)
	menu.SetShowStatusBar(false)
	menu.SetShowTitle(false)
}

func SetListFocus(menu *list.Model, focused bool) {
	delegate := list.NewDefaultDelegate()
	if !focused {
		inactiveColor := lipgloss.Color("#777777")
		delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
			Foreground(inactiveColor).
			BorderForeground(inactiveColor)
		delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
			Foreground(inactiveColor).
			BorderForeground(inactiveColor)
	}

	menu.SetDelegate(delegate)
}
