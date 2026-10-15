package quest

import (
	"solopg/app/cmdrun/domain/card"
)

type Quest struct {
	card.Card


	Done bool
	Prioty Priority
	Rank Rank
	Reward string
}

type Rank string
var (
	// Ou a, b, c , s, ss ..
	MainKind Rank = "main"
	SideKind Rank = "side"
)
// Yaml ca en data drive
type Priority string 
var (
	LowPriority Priority = "low"
	MediumPriority Priority = "medium"
	HighPriority Priority = "high"
)