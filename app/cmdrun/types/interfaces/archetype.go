package interfaces

import "solopg/app/cmdrun/domain/card/attributes/stats"

type Archetype interface {
	GetName() string
	IsPlayable() bool
	GetBonus() []stats.Modifier
}
