package style

import (
	"scne/config"

	"charm.land/lipgloss/v2"
)

var theme = config.GetConfig().Theme

// Palette shared by every style, set by the theme block of the config; never inline a color.
var (
	ColorAccent  = lipgloss.Color(theme.Accent)
	ColorText    = lipgloss.Color(theme.Text)
	ColorMuted   = lipgloss.Color(theme.Muted)
	ColorSubtle  = lipgloss.Color(theme.Subtle)
	ColorError   = lipgloss.Color(theme.Error)
	ColorSuccess = lipgloss.Color(theme.Success)
)
