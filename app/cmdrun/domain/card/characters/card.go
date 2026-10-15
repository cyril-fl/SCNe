package characters

import (
	"solopg/app/cmdrun/domain/card"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/attributes/variety"
	"solopg/app/cmdrun/domain/card/characters/classes"
	"solopg/app/cmdrun/domain/card/characters/races"
	"solopg/app/cmdrun/domain/card/characters/wallet"
	"solopg/app/cmdrun/domain/card/objects"
	"solopg/app/cmdrun/domain/card/objects/equipment"
	"solopg/app/shared/services/logs"
)

type Character struct {
	card.Card

	Class string
	Race  string
	Stats stats.Stats

	Equipment equipment.Equipment
	Inventory []objects.Object
	Wallet    wallet.Wallet

	/*
		LOW Ajouter des "effet" au personnage, ce ne serait pas des effet comme une arme qui fais x ou y mais plus des stats comme un empoissenememnt ex lier a une potion au autre
		ce serais aussi utiliser dans le formulaire de creation de perso pour simulier les effet d'un equipment qu'on a pas ecrit piece par piece.
	*/
}

type Template struct {
	Name        string
	Description string
	Rarity      rarity.Rarity
	Variety     variety.Variety
	Class       string
	Race        string
	Stats       stats.Stats
	Wallet      wallet.Wallet
	Equipment   equipment.Equipment
	Inventory   []objects.Object
}

func New(params Template) (*Character, error) {
	if !params.Variety.Validate() {
		params.Variety = variety.MakeDefault("character_card")
	}

	if !params.Rarity.Validate() {
		params.Rarity = rarity.MakeDefault("")
	}

	newCard, err := card.NewCard(card.NewCardParams{
		Name:        params.Name,
		Description: params.Description,
		Rarity:      params.Rarity,
		Variety:     params.Variety,
	})

	if newCard == nil {
		return nil, logs.Error("error.invalid:new", map[string]any{
			"Subject": params.Name,
			"Error":   err,
		})
	}

	if !classes.Assert(params.Class) {
		return nil, logs.Error("error.invalid", map[string]any{
			"Subject": "class",
			"Receive": params.Class,
		})
	}

	if !races.Assert(params.Race) {
		return nil, logs.Error("error.invalid", map[string]any{
			"Subject": "race",
			"Receive": params.Race,
		})
	}

	return &Character{
		Card:      *newCard,
		Class:     params.Class,
		Race:      params.Race,
		Stats:     params.Stats,
		Equipment: params.Equipment,
		Inventory: params.Inventory,
		Wallet:    params.Wallet,
	}, nil
}

func (character *Character) EquipEquipement(gears []equipment.Gear) {
	for _, gear := range gears {
		character.Equipment.AddGear(gear)
	}
}

func (character *Character) EquipGear(gear equipment.Gear) {
	character.Equipment.AddGear(gear)
}
