package onboardforgecharacter

import (
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/services/factory/millstats"
)

func drowBuild() ([]stats.Modifier, error) {
	generator := millstats.New()
	generator.GenerateFromOracle()

	if generator.HasErr() {
		return nil, generator.GetErr()
	}

	return generator.GetModifiers(), nil
}
