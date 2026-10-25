package style

import "charm.land/lipgloss/v2"

var Panel = lipgloss.NewStyle().
	Padding(0, 1).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorSubtle)
