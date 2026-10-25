package style

import "charm.land/lipgloss/v2"

func Title() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorText())
}

// Collapsed is the title of an unfocused navigation section.
func Collapsed() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorMuted())
}
