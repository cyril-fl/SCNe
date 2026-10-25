package locations

import (
	"scne/app/cmdrun/domain/card"
	"scne/app/cmdrun/domain/card/attributes/rarity"
	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/attributes/variety"
	"scne/app/shared/services/logs"
)

// Location is a location card with its effects.
type Location struct {
	card.Card

	Effects []stats.Effect
}

// Template configures a Location.
type Template struct {
	Name        string
	Description string
	Rarity      rarity.Rarity
	Variety     variety.Variety
	Effects     []stats.Effect
}

// New returns a Location, or an error when the card is invalid.
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
