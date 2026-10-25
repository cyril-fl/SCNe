package style

import "charm.land/lipgloss/v2"

// Key renders a shortcut command in the footer.
var Key = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorText)
