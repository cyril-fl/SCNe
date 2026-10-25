package style

import "charm.land/lipgloss/v2"

func Muted() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorMuted())
}
