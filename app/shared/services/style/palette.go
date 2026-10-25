package style

import (
	"image/color"

	"scne/app/shared/services/theme"

	"charm.land/lipgloss/v2"
)

// Palette shared by every style, read from the current theme at call time; never inline a color.

func ColorAccent() color.Color  { return lipgloss.Color(theme.Get().Accent) }
func ColorText() color.Color    { return lipgloss.Color(theme.Get().Text) }
func ColorMuted() color.Color   { return lipgloss.Color(theme.Get().Muted) }
func ColorSubtle() color.Color  { return lipgloss.Color(theme.Get().Subtle) }
func ColorError() color.Color   { return lipgloss.Color(theme.Get().Error) }
func ColorSuccess() color.Color { return lipgloss.Color(theme.Get().Success) }
