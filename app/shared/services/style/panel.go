package style

import "charm.land/lipgloss/v2"

func Panel() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSubtle())
}
