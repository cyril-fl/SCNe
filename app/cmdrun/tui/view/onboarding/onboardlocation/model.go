package onboardlocation

import (
	"scne/app/cmdrun/domain/card/locations"
	"scne/app/cmdrun/domain/gameplay/portal"
	"scne/app/shared/components/models"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

type model struct {
	reroll models.RerollModel[*locations.Location]
}

func NewModel() *model {
	return &model{
		reroll: models.NewRerollModel(
			models.NewOptionsModel(models.DefaultRerollOptions),
			portal.Teleport,
			3,
		),
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch plumtypes.Key(msg.String()) {
		case plumtypes.KEY_ENTER:
			return m, m.reroll.HandleEnterInput()
		case plumtypes.KEY_ESC:
			// NOTE: Intercept plumtypes.KEY_ESC to prevent list.Update() from quitting.
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.reroll.Choice, cmd = m.reroll.Choice.Update(msg)
	return m, cmd
}

func (m model) HandlesEscape() bool {
	return true
}
