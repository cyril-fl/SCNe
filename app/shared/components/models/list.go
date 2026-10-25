package models

import (
	"charm.land/bubbles/v2/list"
	"solopg/app/shared/services/style"
)

// Item is a list entry carrying a typed value.
type Item[T any] struct {
	title       string
	description string
	value       T
}

// NewItem returns an Item; title is also its filter value.
func NewItem[T any](title, description string, value T) Item[T] {
	return Item[T]{
		title:       title,
		description: description,
		value:       value,
	}
}

// Title returns the item title.
func (i Item[T]) Title() string { return i.title }

// Description returns the item description.
func (i Item[T]) Description() string { return i.description }

// FilterValue returns the title.
func (i Item[T]) FilterValue() string { return i.title }

// Value returns the carried value.
func (i Item[T]) Value() T { return i.value }

// ConfigureList hides the title, filter, pagination, help and status bar of menu.
func ConfigureList(menu *list.Model) {
	menu.Title = ""
	menu.SetShowFilter(false)
	menu.SetShowPagination(false)
	menu.SetShowHelp(false)
	menu.SetShowStatusBar(false)
	menu.SetShowTitle(false)
}

// SetListFocus styles the selected item as active or inactive.
func SetListFocus(menu *list.Model, focused bool) {
	delegate := list.NewDefaultDelegate()
	if !focused {
		inactiveColor := style.Inactive.GetForeground()
		delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
			Foreground(inactiveColor).
			BorderForeground(inactiveColor)
		delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
			Foreground(inactiveColor).
			BorderForeground(inactiveColor)
	}

	menu.SetDelegate(delegate)
}

// ListContentHeight returns the height needed to show every visible item, one item when empty.
func ListContentHeight(menu *list.Model) int {
	delegate := list.NewDefaultDelegate()
	itemCount := len(menu.VisibleItems())

	if itemCount == 0 {
		return delegate.Height()
	}

	return itemCount * (delegate.Height() + delegate.Spacing())
}
