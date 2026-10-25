package inittheme

import (
	"scne/app/shared/services/theme"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type run struct {
	plumprimitive.Process

	cfg theme.Config
}

// Process returns a process overriding the default theme with cfg.
func Process(cfg theme.Config) *run {
	return &run{
		cfg: cfg,
	}
}

// Run installs cfg as the current theme, the default color filling each empty one.
func (p *run) Run() {
	theme.SetCache(p.cfg.Merge(theme.Default))
}

// Getters & Setters

// Methods

// Helpers
