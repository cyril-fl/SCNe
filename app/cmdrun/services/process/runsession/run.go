package cmdrunrunsession

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/characters/classes"
	"solopg/app/cmdrun/domain/card/characters/races"
	"solopg/app/cmdrun/domain/card/locations"
	"solopg/app/cmdrun/services/process/generatesession"
	cmdruntui "solopg/app/cmdrun/tui"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/mango"
	"solopg/app/shared/services/mango/repository"
	"solopg/app/shared/services/process/generatesteplist"

	"solopg/app/cmdrun/services/process/savestate"
	"solopg/app/cmdrun/tui/view/gameboard"
	"solopg/app/cmdrun/tui/view/loadsave"
	"solopg/app/cmdrun/tui/view/onboarding/onboardarchetype"
	"solopg/app/cmdrun/tui/view/onboarding/onboardename"
	"solopg/app/cmdrun/tui/view/onboarding/onboardforgecharacter"
	"solopg/app/cmdrun/tui/view/onboarding/onboardlocation"
	"solopg/app/shared/components/models"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type run struct {
	plumprimitive.Process
	cache

	db  *mango.Mongo
	ctx *cmdruntui.Context
}

type cache struct {
	view *generatesteplist.Generator[cmdruntui.Context]
}

// Process returns the process running the game TUI.
func Process(db *mango.Mongo) *run {
	return &run{
		db: db,
	}
}

// Run builds the save selection and onboarding steps and runs the TUI.
func (p *run) Run() {
	p.makeTuiView()
	p.makeSelectsaveStep()
	p.makeResolveStep()
	p.runTuiView()
}

// Getters & Setters

// GetResult is not implemented; it only logs a warning.
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
		models.Bind(loadsave.NewModel(campaigns), "campaign",
			func(ctx *cmdruntui.Context, selected *campaign.Campaign) {
				ctx.SelectedSave = selected
				if ctx.SelectedSave == nil {
					onboardingSteps := makeOnboardingSteps()
					p.cache.view.Insert(onboardingSteps...)
				}
			}))
}

func (p *run) makeResolveStep() {
	if p.HasErr() {
		return
	}

	p.cache.view.Add(
		models.Step[cmdruntui.Context]{
			Submodel: models.Resolve(),
			Resolve: func(ctx *cmdruntui.Context, value any) error {
				process := generatesession.Process(p.db, ctx)
				process.Run()

				if process.HasErr() {
					return logs.Error("error.unexpected:action", map[string]any{
						"Action": "unexpected:action.build_session",
						"Error":  process.GetErr().Error(),
					})
				}

				p.cache.view.Insert(models.Step[cmdruntui.Context]{
					Submodel: gameboard.NewModel(gameboard.UiParams{
						Session: process.GetResult(),
						OnSave:  savestate.Save(p.db, process.GetResult()),
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
func makeOnboardingSteps() []models.Step[cmdruntui.Context] {
	return []models.Step[cmdruntui.Context]{
		models.Bind(onboardename.NewModel(), "name",
			func(ctx *cmdruntui.Context, name string) {
				ctx.SelectedName = name

				logs.SilentInfo("onboarding.selected", map[string]any{
					"Subject": "name",
					"Value":   name,
				})
			}),
		models.Bind(onboardarchetype.NewModel("races", races.List()), "race",
			func(ctx *cmdruntui.Context, selected models.Item[races.Race]) {
				race := selected.Value()
				ctx.SelectedRace = &race

				logs.SilentInfo("onboarding.selected", map[string]any{
					"Subject": "race",
					"Value":   race.GetName(),
				})
			}),
		models.Bind(onboardarchetype.NewModel("classes", classes.List()), "class",
			func(ctx *cmdruntui.Context, selected models.Item[classes.Class]) {
				class := selected.Value()
				ctx.SelectedClass = &class

				logs.SilentInfo("onboarding.selected", map[string]any{
					"Subject": "class",
					"Value":   class.GetName(),
				})
			}),
		models.Bind(onboardforgecharacter.NewModel(), "build",
			func(ctx *cmdruntui.Context, build []stats.Modifier) {
				ctx.SelectedBuild = build

				logs.SilentInfo("onboarding.selected", map[string]any{
					"Subject": "build",
					"Value":   build,
				})
			}),
		models.Bind(onboardlocation.NewModel(), "location",
			func(ctx *cmdruntui.Context, location *locations.Location) {
				ctx.SelectedLocation = location

				logs.SilentInfo("onboarding.selected", map[string]any{
					"Subject": "location",
					"Value":   location,
				})
			}),
	}
}
