package quest

import (
	"solopg/app/cmdrun/domain/card"
	"solopg/app/cmdrun/domain/card/attributes/priority"
	"solopg/app/cmdrun/domain/card/attributes/rank"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/shared/services/logs"
)

type Quest struct {
	card.Card

	Done     bool
	Priority priority.Priority
	Rank     rank.Rank
	Reward   string
}

type Template struct {
	Name        string
	Description string
	Variety     variety.Variety

	Done     bool
	Priority priority.Priority
	Rank     rank.Rank
	Reward   string
}

func New(params Template) (*Quest, error) {
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
		Card:     *newCard,
		Done:     params.Done,
		Priority: params.Priority,
		Rank:     params.Rank,
		Reward:   params.Reward,
	}, nil
}
