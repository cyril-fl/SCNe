package onboardforgecharacter

import (
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/shared/components/models"

	sharedtui "solopg/app/shared/tui"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	reroll models.RerollModel[[]stats.Modifier]
}

func NewModel() *model {
	return &model{
		reroll: models.NewRerollModel(
			models.NewOptionsModel(models.DefaultRerollOptions),
			drowBuild,
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
		switch msg.String() {
		case sharedtui.KEY_ENTER:
			return m, m.reroll.HandleEnterInput()
		}
	}

	var cmd tea.Cmd
	m.reroll.Choice, cmd = m.reroll.Choice.Update(msg)
	return m, cmd
}
