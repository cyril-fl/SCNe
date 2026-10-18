package resolvecampaign

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/cmdrun/domain/card/attributes/rarity"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/characters"
	"solopg/app/cmdrun/domain/card/characters/wallet"
	"solopg/app/cmdrun/domain/card/objects"
	"solopg/app/cmdrun/domain/card/objects/equipment"
	cmdruntui "solopg/app/cmdrun/tui"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"strings"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type run struct {
	plumprimitive.Process
	cache

	ctx *cmdruntui.Context
}

type cache struct {
	campaign *campaign.Campaign
	player   *characters.Character

	err error
}

func Process(ctx *cmdruntui.Context) *run {
	return &run{
		ctx: ctx,
	}
}

func (p *run) Run() {
	if p.isNewCampaign() {
		p.assertContext()
		p.generateCharacter()
		p.generateCampaign()
	} else {
		p.cache.campaign = p.ctx.SelectedSave
	}
}

// Getters & Setters
func (p *run) GetResult() *campaign.Campaign {
	return p.cache.campaign
}

// Methods
func (p *run) isNewCampaign() bool {
	return p.ctx.SelectedSave == nil
}

func (p *run) assertContext() {
	isValid := false
	c := p.ctx

	for _, v := range []any{c.SelectedRace, c.SelectedClass, c.SelectedLocation} {
		if v == nil {
			continue
		}
		isValid = true
		break
	}

	if isValid {
		return
	}

	p.cache.err = logs.Error("error.required", map[string]any{
		"Subject":  "campaign",
		"Property": strings.Join([]string{i19n.Localize("race"), i19n.Localize("class"), i19n.Localize("location")}, ", "),
	})

	p.checkError()
}

func (p *run) generateCharacter() {
	if p.HasErr() {
		return
	}

	p.cache.player, p.cache.err = characters.New(characters.Template{
		Name:      p.ctx.SelectedName,
		Race:      p.ctx.SelectedRace.GetName(),
		Class:     p.ctx.SelectedClass.GetName(),
		Rarity:    rarity.Default(),
		Stats:     makeStatsFromContext(p.ctx),
		Equipment: makeEquipementFromContext(p.ctx),
		Inventory: []objects.Object{},
		Wallet:    wallet.Wallet{},
	})

	p.checkError()
}

func (p *run) generateCampaign() {
	if p.HasErr() {
		return
	}

	p.cache.campaign = campaign.New(campaign.Template{
		Player:          p.cache.player,
		CurrentLocation: p.ctx.SelectedLocation,
	})
}

func (p *run) checkError() {
	if p.cache.err != nil {
		p.SetErr(p.cache.err)
		p.cache.err = nil
	}
}

// Helpers
func makeEquipementFromContext(ctx *cmdruntui.Context) equipment.Equipment {
	name := ctx.SelectedClass.GetEquipementName()
	set := equipment.FindEquipementByName(name)

	return equipment.NewSet(set)
}

func makeStatsFromContext(ctx *cmdruntui.Context) stats.Stats {
	modifiers := makeModifiersFromContext(ctx)

	baseStats := stats.GetBasic()
	baseStats.ApplyModifiers(modifiers)

	return baseStats
}

func makeModifiersFromContext(ctx *cmdruntui.Context) []stats.Modifier {
	raceBoost := ctx.SelectedRace.GetBonus()
	classBoost := ctx.SelectedClass.GetBonus()
	build := ctx.SelectedBuild

	modifiers := make([]stats.Modifier, 0, len(raceBoost)+len(classBoost)+len(build))
	modifiers = append(modifiers, raceBoost...)
	modifiers = append(modifiers, classBoost...)
	modifiers = append(modifiers, build...)

	return modifiers
}
