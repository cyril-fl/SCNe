package cmdservetui

import (
	"solopg/app/shared/components/models"
	"solopg/app/shared/messages"

	plumtypes "github.com/cyril-fl/plum/types"

	tea "charm.land/bubbletea/v2"
)

// Context is the empty step context of the serve TUI.
type Context struct{}

type model struct {
	steps models.StepModel[Context]
}

// NewModel returns the serve TUI running steps.
func NewModel(steps *models.ContextStepList[Context]) model {
	return model{
		steps: models.NewStepModel(steps),
	}
}

// Init starts the current step.
func (m *model) Init() tea.Cmd {
	if submodel := m.steps.GetCurrentSubmodel(); submodel != nil {
		return submodel.Init()
	}

	return nil
}

// Update quits on ctrl+c or ctrl+q and forwards other messages to the current step.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.steps.SetSize(&msg)
		return m, m.steps.UpdateSubmodel(msg)
	case tea.KeyMsg:
		switch plumtypes.Key(msg.String()) {
		case plumtypes.SHORT_CTRL_C, plumtypes.SHORT_CTRL_Q:
			return m, tea.Quit
		}
		return m, m.steps.UpdateSubmodel(msg)

	case messages.Resolved:
		return m, m.steps.ResolveSubmodel(msg)

	default:
		return m, m.steps.UpdateSubmodel(msg)
	}
}
