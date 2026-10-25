package style

import "charm.land/lipgloss/v2"

func Error() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorError())
}
