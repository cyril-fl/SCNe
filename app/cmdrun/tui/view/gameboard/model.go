package gameboard

import (
	"solopg/app/cmdrun/services/factory/roll"
	"solopg/app/cmdrun/services/game"
	"solopg/app/cmdrun/tui/view/sidemenu/codexmenu"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/chat"
	"solopg/app/shared/components/formpage"
	"solopg/app/shared/components/navigation"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/types/primitive"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

// - Model - //
type model struct {
	primitive.Fallible

	state  statemachine
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
		state: newstate(params),
		menu: navigation.New(navigation.Template{
			Static: []string{makeStatView(params.Engine)},
			Parts:  makePanelPart(params.Engine),
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
		m.state.HandleTransition(CHAT_SECTION, state(msg.ID))
		return m, nil
	case formpage.CloseMsg:
		m.state.HandleTransition(m.state.Current(), CHAT_SECTION)
		return m, nil
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	case roll.RollResultMsg:
		m.state.Amend(msg.Content)
		return m, nil

	// Cmd messages
	// case codexmenu.Msg:
	// 	return m.handleCodexAction(msg)
	// case dicemenu.Msg:
	// 	return m.handleDiceRolled(msg)
	// case hintmenu.Msg:
	// 	return m.handleHintRolled(msg)
	// case oraclemenu.Msg:
	// 	return m.handleOracleRolled(msg)

	case sharedtui.SaveMsg:
		return m.handleSaveInput(msg)
	// case sharedtui.ScrollMsg:
	// 	return m.handleViewportScroll(msg)
	// case sharedtui.Refresh:
	// 	return m.handleViewRefresh(msg)
	case sharedtui.ErrorMsg:
		return m.handleError(msg)

	}

	return m, nil
}

// Helpers
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
