package style

import "charm.land/lipgloss/v2"

func Accent() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorAccent())
}
