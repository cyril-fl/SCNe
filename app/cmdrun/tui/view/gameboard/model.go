package gameboard

import (
	"solopg/app/cmdrun/services/game"
	"solopg/app/cmdrun/tui/view/sidemenu/codexmenu"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/chat"
	"solopg/app/shared/components/formpage"
	"solopg/app/shared/components/navigation"
	"solopg/app/shared/messages"
	"solopg/app/shared/services/statemachine"
	"solopg/app/shared/types/primitive"
	"time"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

// - Model - //
type model struct {
	primitive.Fallible

	state  statemachine.ViewState[state]
	menu   *navigation.Navigation
	engine *game.Engine
	save   func() error
	size   size.Size
}

type UiParams struct {
	Engine *game.Engine
	OnSave func() error
}

func NewModel(params UiParams) *model {
	return &model{
		state: makeViewState(params),
		menu: navigation.New(navigation.Template{
			Static:     []string{makeStatView(params.Engine)},
			Scrollable: makePanelPart(params.Engine),
		}),
		engine: params.Engine,
		save:   params.OnSave,
	}
}

// Getters & Setters

// Methods
func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	statecmd := m.state.Update(msg)
	var menucmd tea.Cmd
	if m.state.Current() == CHAT_SECTION {
		menucmd = m.menu.HandleUpdate(msg)
	}
	next, cmd := m.handleUpdate(msg)

	return next, tea.Sequence(menucmd, statecmd, cmd)
}

func (m *model) handleUpdate(msg tea.Msg) (*model, tea.Cmd) {
	switch msg := msg.(type) {
	case codexmenu.OpenPageMsg:
		m.state.HandleTransition(state(msg.ID))
	case formpage.CloseMsg:
		m.state.HandleTransition(CHAT_SECTION)
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	case messages.Result:
		return m, tea.Tick(1500*time.Millisecond, func(time.Time) tea.Msg {
			return messages.Delayed{Content: msg.Content}
		})
	case messages.Delayed:
		m.state.Amend(msg.Content)
	case messages.Saved:
		return m.handleSaveInput(msg)
	case messages.Failed:
		return m.handleError(msg)
	}

	return m, nil
}

func makeChat(engine *game.Engine) *chat.Chat {
	return chat.New(chat.Template{
		Author:  initAuthor(engine),
		History: initJournal(engine),
		Register: func(author, message string) {
			engine.GetAdventureJournal().AddEntry(author, message)
		},
	})
}

func initAuthor(engine *game.Engine) string {
	author := "System"
	if engine.State.Player != nil {
		author = engine.State.Player.Name
	}
	return author
}

func initJournal(engine *game.Engine) []string {
	raw := make([]string, 0, len(engine.State.AdventureLog.Entries))
	for _, entry := range engine.State.AdventureLog.Entries {
		raw = append(raw, entry.String())
	}
	return raw
}

// HandlesEscape lets the global router forward Escape while a Codex page is open.
func (m model) HandlesEscape() bool {
	return true
}

// Helpers
