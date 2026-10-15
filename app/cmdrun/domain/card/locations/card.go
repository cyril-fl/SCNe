package locations

import (
	"solopg/app/cmdrun/domain/card"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/shared/services/logs"
)

type Location struct {
	card.Card

	Effects []stats.Effect
}

type Template struct {
	Name        string
	Description string
	Rarity      rarity.Rarity
	Variety     variety.Variety
	Effects     []stats.Effect
}

func New(params Template) (*Location, error) {
	newCard, err := card.NewCard(card.NewCardParams{
		Name:        params.Name,
		Description: params.Description,
		Rarity:      params.Rarity,
		Variety:     params.Variety,
	})

	if newCard == nil {
		return nil, logs.Error("error.invalid:new", map[string]any{
			"Subject": "location",
			"Error":   err,
		})
	}

	return &Location{
		Card:    *newCard,
		Effects: params.Effects,
	}, nil
}