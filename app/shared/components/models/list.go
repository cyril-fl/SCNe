package models

import (
	"scne/app/shared/services/style"

	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
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

// NewList returns a focused, configured list whose delegate fits items.
func NewList(items []list.Item, width, height int) list.Model {
	model := list.New(items, newDelegate(items, true), width, height)
	ConfigureList(&model)
	return model
}

// SetListFocus styles the selected item as active or inactive.
// The delegate is rebuilt from the current items, so call it again after SetItems.
func SetListFocus(menu *list.Model, focused bool) {
	menu.SetDelegate(newDelegate(menu.Items(), focused))
}

// ListContentHeight returns the height needed to show every visible item, one item when empty.
func ListContentHeight(menu *list.Model) int {
	delegate := newDelegate(menu.Items(), false)
	itemCount := len(menu.VisibleItems())

	if itemCount == 0 {
		return delegate.Height()
	}

	return itemCount * (delegate.Height() + delegate.Spacing())
}

// newDelegate shows descriptions, and spaces items, only when one item has a description.
func newDelegate(items []list.Item, focused bool) list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = hasDescription(items)
	if !delegate.ShowDescription {
		delegate.SetHeight(1)
		delegate.SetSpacing(0)
	}

	marker, selected := style.ColorAccent, style.ColorText
	if !focused {
		marker, selected = style.ColorMuted, style.ColorMuted
	}

	normal := lipgloss.NewStyle().Foreground(style.ColorText).PaddingLeft(2)
	chevron := lipgloss.NewStyle().
		Border(lipgloss.Border{Left: "›"}, false, false, false, true).
		BorderForeground(marker).
		PaddingLeft(1)

	delegate.Styles.NormalTitle = normal
	delegate.Styles.NormalDesc = normal.Foreground(style.ColorMuted)
	delegate.Styles.SelectedTitle = chevron.Foreground(selected).Bold(focused)
	delegate.Styles.SelectedDesc = normal.Foreground(style.ColorMuted)
	delegate.Styles.DimmedTitle = normal.Foreground(style.ColorMuted)
	delegate.Styles.DimmedDesc = normal.Foreground(style.ColorSubtle)

	return delegate
}

func hasDescription(items []list.Item) bool {
	for _, item := range items {
		if item, ok := item.(list.DefaultItem); ok && item.Description() != "" {
			return true
		}
	}
	return false
}
