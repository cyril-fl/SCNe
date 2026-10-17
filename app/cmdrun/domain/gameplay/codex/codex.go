package codex

import (
	"solopg/app/cmdrun/domain/gameplay/codex/service/factory/milltable"
	codextable "solopg/app/cmdrun/domain/gameplay/codex/table"
)

type Codex struct {
	Npcs      codextable.Table
	Beastiary codextable.Table
	Locations codextable.Table
	Objects   codextable.Table
	Quest     codextable.Table
}

func New() *Codex {
	codex := &Codex{}
	codex.EnsureInitialized()
	return codex
}

// EnsureInitialized repairs a Codex loaded from storage where nested tables may be nil.
func (c *Codex) EnsureInitialized() {
	mills := map[string]milltable.MillTable{
		"npcs":      milltable.Npcs(&c.Npcs),
		"beastiary": milltable.Beastiary(&c.Beastiary),
		"locations": milltable.Locations(&c.Locations),
		"objects":   milltable.Objects(&c.Objects),
		"quest":     milltable.Quest(&c.Quest),
	}

	for _, mill := range mills {
		mill.Make()
	}
}
