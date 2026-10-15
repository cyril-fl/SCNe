package card

import (
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/types/primitive"
	"strings"
)

type Card struct {
	primitive.Descriptible

	Rarity  rarity.Rarity
	Variety variety.Variety
}

type NewCardParams struct {
	Name        string
	Description string
	Rarity      rarity.Rarity
	Variety     variety.Variety
}

func NewCard(params NewCardParams) (*Card, error) {
	if !params.Rarity.Validate() {

		valid := make([]string, 0, len(rarity.List()))
		for _, r := range rarity.List() {
			valid = append(valid, i19n.Localize(r.String()))
		}

		return nil, logs.Error("error.unexpected:value", map[string]any{
			"Subject":  "rarity",
			"Expected": strings.Join(valid, ", "),
			"Value":    params.Rarity,
		})
	}

	if !params.Variety.Validate() {

		valid := make([]string, 0, len(rarity.List()))
		for _, v := range variety.List() {
			valid = append(valid, i19n.Localize(v.String()))
		}

		return nil, logs.Error("error.unexpected:value", map[string]any{
			"Subject":  "variety",
			"Expected": strings.Join(valid, ", "),
			"Value":    params.Variety,
		})
	}

	return &Card{
		Descriptible: primitive.Describe(primitive.TemplateDescriptible{
			Name:        params.Name,
			Description: params.Description,
		}),
		Rarity:  params.Rarity,
		Variety: params.Variety,
	}, nil
}
