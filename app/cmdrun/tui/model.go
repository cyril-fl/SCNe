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
		steps:  models.NewStepModel(steps),
	}
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

// TODO voitr si c'est vraiment a la bonne place ici
type DelegateUpdateFunc func(msg tea.Msg) (tea.Model, tea.Cmd)

type EscapeSupport interface {
	HandlesEscape() bool
}
