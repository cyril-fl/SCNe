package onboardforgecharacter

import (
	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/services/factory/millstats"
)

func drowBuild() ([]stats.Modifier, error) {
	generator := millstats.New()
	generator.GenerateFromOracle()

	if generator.HasErr() {
		return nil, generator.GetErr()
	}

	return generator.GetModifiers(), nil
}
