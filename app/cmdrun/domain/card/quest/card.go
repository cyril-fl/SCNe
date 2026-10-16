package quest

import (
	"solopg/app/cmdrun/domain/card"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/shared/services/logs"
)

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

type Quest struct {
	card.Card

	Done bool
	Prioty Priority
	Rank Rank
	Reward string
}

type Template struct {
	Name        string
	Description string
	Variety     variety.Variety

	Done bool
	Prioty Priority
	Rank Rank
	Reward string
}

func New(params Template)(*Quest, error){
	// CHECK cpparce que nil ca a l'ai de passer
	if !params.Variety.Validate() {
		params.Variety = variety.MakeDefault("quest_card")
	}

	newCard, err := card.NewCard(card.NewCardParams{
		Name:        params.Name,
		Description: params.Description,
		Rarity:      rarity.MakeDefault(""),
		Variety:     params.Variety,
	})

	if newCard == nil {
		return nil, logs.Error("error.invalid:new", map[string]any{
			"Subject": "item",
			"Error":   err,
		})
	}

	return &Quest{
		Name: params.Name,
		Description: params.Description,
		Done: params.Done,
		Prioty: params.Prioty,
		Rank: params.Rank,
		Reward: params.Reward,
	}, nil
}