package cmdrunrunsession

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/characters/classes"
	"solopg/app/cmdrun/domain/card/characters/races"
	"solopg/app/cmdrun/domain/card/locations"
	"solopg/app/cmdrun/services/process/generatengine"
	cmdruntui "solopg/app/cmdrun/tui"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango"
	"solopg/app/shared/services/mango/repository"
	"solopg/app/shared/services/process"
	"solopg/app/shared/services/process/generatesteplist"

	"solopg/app/cmdrun/services/process/savestate"
	"solopg/app/cmdrun/tui/models"
	"solopg/app/cmdrun/tui/view/gameboard"
	"solopg/app/cmdrun/tui/view/loadsave"
	"solopg/app/cmdrun/tui/view/onboarding/onboardarchetype"
	"solopg/app/cmdrun/tui/view/onboarding/onboardename"
	"solopg/app/cmdrun/tui/view/onboarding/onboardforgecharacter"
	"solopg/app/cmdrun/tui/view/onboarding/onboardlocation"
)

type run struct {
	process.Process
	cache

	db  *mango.Mongo
	ctx *cmdruntui.Context
}

type cache struct {
	view *generatesteplist.Generator[cmdruntui.Context]
}

func Process(db *mango.Mongo) *run {
	return &run{
		db: db,
	}
}

func (p *run) Run() {
	p.makeTuiView()
	p.makeSelectsaveStep()
	p.makeResolveStep()
	p.runTuiView()
}

// Getters & Setters
func (p *run) GetResult() {
	logs.SilentWarning("error:not_implemented", map[string]any{
		"Function": "GetResult",
		"Subject":  "session",
	})
}

// Methods
func (p *run) makeTuiView() {
	p.cache.view = generatesteplist.Process[cmdruntui.Context]()
}

func (p *run) makeSelectsaveStep() {
	if p.HasErr() {
		return
	}

	repository := repository.Campaign().SetDb(p.db)
	campaigns, err := repository.Load()
	if err != nil {
		err := logs.Error("error.unexpected:action", map[string]any{
			"Action": "unexpected:action.load_campaigns",
			"Value":  err.Error(),
		})

		p.SetErr(err)

		return
	}

	p.cache.view.Add(
		cmdruntui.Step{
			Submodel: loadsave.NewModel(campaigns),
			Resolve: func(ctx *cmdruntui.Context, value any) error {

				selected, ok := value.(*campaign.Campaign)
				if !ok {
					return logs.Error("error.unexpected:value", map[string]any{
						"Subject":  "campaign",
						"Expected": "*campaign.Campaign",
						"Value":    value,
					})
				}

				ctx.SelectedSave = selected
				if ctx.SelectedSave == nil {
					onboardingSteps := makeOnboardingSteps()
					p.cache.view.Insert(onboardingSteps...)
				}

				return nil
			},
		})
}

func (p *run) makeResolveStep() {
	if p.HasErr() {
		return
	}

	p.cache.view.Add(
		cmdruntui.Step{
			Submodel: models.Resolve(),
			Resolve: func(ctx *cmdruntui.Context, value any) error {
				process := generatengine.Process(p.db, ctx)
				process.Run()

				if process.HasErr() {
					return logs.Error("error.unexpected:action", map[string]any{
						"Action": "unexpected:action.build_engine",
						"Error":  process.GetErr().Error(),
					})
				}

				p.cache.view.Insert(cmdruntui.Step{
					Submodel: gameboard.NewModel(gameboard.UiParams{
						Engine: process.GetResult(),
						OnSave: savestate.Save(p.db, process.GetResult()),
					}),
				})

				return nil
			},
		})
}

func (p *run) runTuiView() {
	model := cmdruntui.NewModel(p.cache.view.GetSteps())

	p.cache.view.Run(&model)

	if p.cache.view.HasErr() {
		p.SetErr(p.cache.view.GetErr())
	}
}

// Helpers
func makeOnboardingSteps() []cmdruntui.Step {
	return []cmdruntui.Step{
		{
			Submodel: onboardename.NewModel(),
			Resolve: func(ctx *cmdruntui.Context, value any) error {
				name, ok := value.(string)
				if !ok {
					return logs.Error("error.unexpected:value", map[string]any{
						"Subject":  "name",
						"Expected": "string",
						"Value":    value,
					})
				}
				ctx.SelectedName = name

				logs.SilentInfo("onboarding.selected", map[string]any{
					"Subject": "name",
					"Value":   name,
				})

				return nil
			},
		},
		{
			Submodel: onboardarchetype.NewModel("races", races.List()),
			Resolve: func(ctx *cmdruntui.Context, value any) error {
				selected, ok := value.(models.Item[races.Race])
				if !ok {
					return logs.Error("error.unexpected:value", map[string]any{
						"Subject":  "race",
						"Expected": "models.Item[races.Race]",
						"Value":    value,
					})
				}
				race := selected.Value()
				ctx.SelectedRace = &race

				logs.SilentInfo("onboarding.selected", map[string]any{
					"Subject": "race",
					"Value":   race.GetName(),
				})

				return nil
			},
		},
		{
			Submodel: onboardarchetype.NewModel("classes", classes.List()),
			Resolve: func(ctx *cmdruntui.Context, value any) error {
				selected, ok := value.(models.Item[classes.Class])
				if !ok {
					return logs.Error("error.unexpected:value", map[string]any{
						"Subject":  "class",
						"Expected": "models.Item[classes.Class]",
						"Value":    value,
					})
				}
				class := selected.Value()
				ctx.SelectedClass = &class

				logs.SilentInfo("onboarding.selected", map[string]any{
					"Subject": "class",
					"Value":   class.GetName(),
				})

				return nil
			},
		},
		{
			Submodel: onboardforgecharacter.NewModel(),
			Resolve: func(ctx *cmdruntui.Context, value any) error {
				build, ok := value.([]stats.Modifier)
				if !ok {

					return logs.Error("error.unexpected:value", map[string]any{
						"Subject":  "build",
						"Expected": "[]stats.Modifier",
						"Value":    value,
					})
				}
				ctx.SelectedBuild = build

				logs.SilentInfo("onboarding.selected", map[string]any{
					"Subject": "build",
					"Value":   build,
				})

				return nil
			},
		},
		{
			Submodel: onboardlocation.NewModel(),
			Resolve: func(ctx *cmdruntui.Context, value any) error {
				location, ok := value.(*locations.Location)
				if !ok {
					return logs.Error("error.unexpected:value", map[string]any{
						"Subject":  "location",
						"Expected": "*locations.Location",
						"Value":    value,
					})
				}
				ctx.SelectedLocation = location

				logs.SilentInfo("onboarding.selected", map[string]any{
					"Subject": "location",
					"Value":   location,
				})

				return nil
			},
		},
	}
}
