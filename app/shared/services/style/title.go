package style

import "charm.land/lipgloss/v2"

var Title = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorText)

// Collapsed is the title of an unfocused navigation section.
var Collapsed = lipgloss.NewStyle().
	Foreground(ColorMuted)
