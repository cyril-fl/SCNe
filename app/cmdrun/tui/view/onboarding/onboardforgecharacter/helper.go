package onboardforgecharacter

import (
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/services/factory/millstats"
)

// func handleWindowResize(m *model, msg tea.WindowSizeMsg) {
// 	// m.reroll.Options.SetSize(msg.Width, max(0, msg.Height-3))
// }

func drowBuild() ([]stats.Modifier, error) {
	generator := millstats.New()
	generator.GenerateFromOracle()

	if generator.HasErr() {
		return nil, generator.GetErr()
	}

	return generator.GetModifiers(), nil
}
