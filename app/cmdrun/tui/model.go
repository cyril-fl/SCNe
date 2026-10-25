package cmdruntui

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/characters/classes"
	"solopg/app/cmdrun/domain/card/characters/races"
	"solopg/app/cmdrun/domain/card/locations"
	"solopg/app/shared/components/layout"
	"solopg/app/shared/components/models"
	"solopg/app/shared/messages"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	layout  layout.Layout
	steps   models.StepModel[Context]
	context Context
}

// Context gathers the save selection and the onboarding choices.
type Context struct {
	SelectedSave     *campaign.Campaign
	SelectedName     string
	SelectedRace     *races.Race
	SelectedClass    *classes.Class
	SelectedBuild    []stats.Modifier
	SelectedLocation *locations.Location
}

// NewModel returns the game TUI running steps inside the layout.
func NewModel(steps *models.ContextStepList[Context]) model {
	return model{
		layout: layout.New(),
		steps:  models.NewStepModel(steps),
	}
}

// Init starts the current step.
func (m model) Init() tea.Cmd {
	if submodel := m.steps.GetCurrentSubmodel(); submodel != nil {
		return submodel.Init()
	}

	return nil
}

// Update quits on ctrl+q, and on esc unless the step handles it,
// and forwards other messages to the current step.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.steps.SetSize(&msg)
		return m, m.steps.UpdateSubmodel(msg)
	case tea.KeyMsg:
		return m.handleEvent(msg)

	case messages.Resolved:
		return m, m.steps.ResolveSubmodel(msg)

	default:
		return m, m.steps.UpdateSubmodel(msg)
	}
}

// EscapeSupport is implemented by steps that may consume esc instead of quitting.
type EscapeSupport interface {
	HandlesEscape() bool
}
