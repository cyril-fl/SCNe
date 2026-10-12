package gameboard

import (
	"solopg/app/cmdrun/services/game"
	"solopg/app/cmdrun/tui/view/sidemenu"
	"solopg/app/shared/components/chat"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/types/primitive"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

// - Model - //
type model struct {
	cache
	primitive.Fallible

	state statemachine

	engine *game.Engine
	save   func() error

	menu            []sidemenu.MenuItem
	activeMenuIndex int
}

type cache struct {
	size *tea.WindowSizeMsg
}

type UiParams struct {
	Engine *game.Engine
	OnSave func() error
}

func NewModel(params UiParams) *model {
	return &model{
		state: newstate(params),

		engine: params.Engine,
		save:   params.OnSave,

		menu:            initSideMenu(params.Engine),
		activeMenuIndex: 0,
	}
}

// Getters & Setters

// Methods
func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// menu := m.getMenuActiveElement()
	// return menu.HandleUpdate(sidemenu.UpdateParams{
	// 	Model:    m,
	// 	Msg:      msg,
	// 	Delegate: m.handleUpdate,
	// 	// Viewport: m.viewport,
	// })
	statecmd := m.state.Update(msg)
	next, cmd := m.handleUpdate(msg)

	return next, tea.Sequence(statecmd, cmd)
}

func (m *model) handleUpdate(msg tea.Msg) (*model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)

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

	case cursor.BlinkMsg:
		return m.handleCursorBlink(msg)
	}

	return m, nil
}

// Helpers
func makeChat(engine *game.Engine) *chat.Chat {
	return chat.New(chat.Template{
		Author: initAuthor(engine),
		History: initJournal(engine),
		Register: func (author, message string) {
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
