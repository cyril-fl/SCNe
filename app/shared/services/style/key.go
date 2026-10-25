package style

import "charm.land/lipgloss/v2"

// Key renders a shortcut command in the footer.
func Key() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorText())
}
