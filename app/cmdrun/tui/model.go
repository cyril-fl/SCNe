package cmdruntui

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/characters/classes"
	"solopg/app/cmdrun/domain/card/characters/races"
	"solopg/app/cmdrun/domain/card/locations"
	"solopg/app/shared/components/layout"
	"solopg/app/shared/components/models"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/types/primitive"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	cache
	primitive.Fallible
	layout layout.Layout

	steps models.StepModel[Context]
	context Context
}

type cache struct {
	size *tea.WindowSizeMsg
}

type Context struct {
	SelectedSave     *campaign.Campaign
	SelectedName     string
	SelectedRace     *races.Race
	SelectedClass    *classes.Class
	SelectedBuild    []stats.Modifier
	SelectedLocation *locations.Location
}

func NewModel(steps *models.ContextStepList[Context]) model {
	return model{
		layout: layout.New(),
		steps:  models.NewStepModel[Context](steps),
	}
}

func (m *model) resolveCurrentStep(value any) error {
	if current := m.steps.CurrentStep(); current != nil && current.Resolve != nil {
		return current.Resolve(&m.context, value)
	}

	return nil
}

func (m *model) selectNextStep() bool {
	for {
		next := m.steps.NextStep()
		if next == nil {
			return false
		}

		if next.Skip == nil || !next.Skip(&m.context) {
			return true
		}
	}
}

func (m *model) forwardNextStep() (*model, tea.Cmd) {
	if !m.selectNextStep() {
		return m, tea.Quit
	}

	submodel := m.steps.GetCurrentSubmodel()
	if submodel == nil {
		return m, tea.Quit
	}

	initCmd := submodel.Init()
	if m.size == nil {
		return m, initCmd
	}

	return m, tea.Sequence(initCmd, func() tea.Msg {
		return *m.size
	})
}

func (m model) Init() tea.Cmd {
	if submodel := m.steps.GetCurrentSubmodel(); submodel != nil {
		return submodel.Init()
	}

	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = &msg
		return m.handleResize(msg)
	case tea.KeyMsg:
		return m.handleEvent(msg)

	case sharedtui.ResolutionMsg:
		return m.handleResolution(msg)

	default:
		return m.handleSubmodelUpdate(msg)
	}
}

type DelegateUpdateFunc func(msg tea.Msg) (tea.Model, tea.Cmd)

type EscapeSupport interface {
	HandlesEscape() bool
}
