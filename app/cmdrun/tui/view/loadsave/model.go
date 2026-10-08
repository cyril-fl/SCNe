package loadsave

import (
	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/shared/components/page"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	cache

	list     list.Model
	selected *campaign.Campaign

	page *page.Page
}
type cache struct {
	size *tea.WindowSizeMsg
}

func NewModel(saves []campaign.Campaign) *model {
	return &model{
		page: page.NewPage(page.Template{
			Title:    "Load/Save Game",
			Subtitle: "Load or save your game progress",
		}),

		list: makeModel(saves),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)
	case tea.KeyMsg:
		return m.handKeyMsg(msg)
	default:
		return m, nil
	}
}
