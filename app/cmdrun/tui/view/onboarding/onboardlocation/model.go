package onboardlocation

import (
	"solopg/app/cmdrun/domain/card/locations"
	"solopg/app/cmdrun/domain/gameplay/portal"
	"solopg/app/cmdrun/tui/models"
	"solopg/app/shared/components/page"
	sharedtui "solopg/app/shared/tui"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	cache
	page *page.Page

	reroll models.RerollModel[*locations.Location]
}

type cache struct {
	size *tea.WindowSizeMsg
}

func NewModel() *model {
	return &model{
		page: page.NewPage(page.Template{
			Title:    "onboarding.character:create",
			Subtitle: "onboarding.character:name",
		}),

		reroll: models.NewRerollModel(
			models.NewOptionsModel(models.DefaultRerollOptions),
			portal.Teleport,
			3,
		),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	handleWindowResize(&m, msg)

	case tea.KeyMsg:
		switch msg.String() {
		case sharedtui.KEY_ENTER:
			return m, m.reroll.HandleEnterInput()
		}
	}

	var cmd tea.Cmd
	m.reroll.Options, cmd = m.reroll.Options.Update(msg)
	return m, cmd
}
