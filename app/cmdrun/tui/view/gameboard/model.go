package gameboard

import (
	"solopg/app/cmdrun/services/game"
	cmdruntui "solopg/app/cmdrun/tui"
	"solopg/app/cmdrun/tui/view/sidemenu"
	"solopg/app/cmdrun/tui/view/sidemenu/codexmenu"
	"solopg/app/cmdrun/tui/view/sidemenu/dicemenu"
	"solopg/app/cmdrun/tui/view/sidemenu/hintmenu"
	"solopg/app/cmdrun/tui/view/sidemenu/oraclemenu"
	"solopg/app/shared/types/primitive"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

// - Model - //
type model struct {
	cache
	primitive.Fallible

	textarea         textarea.Model
	viewport         viewport.Model
	author           string
	chat             []string
	adventurejournal []string

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

func NewModel(params UiParams) model {
	return model{
		viewport:         initViewport(),
		textarea:         initTextarea(),
		author:           initAuthor(params.Engine),
		chat:             make([]string, 0),
		adventurejournal: initJournal(params.Engine),

		engine: params.Engine,
		save:   params.OnSave,

		menu:            initSideMenu(params.Engine),
		activeMenuIndex: 0,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	menu := m.getMenuActiveElement()

	return menu.HandleUpdate(sidemenu.UpdateParams{
		Model:    m,
		Msg:      msg,
		Delegate: m.handleUpdate,
		Viewport: m.viewport,
	})
}

func (m model) handleUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Tea messages
	case tea.WindowSizeMsg:
		return m.handleWindowResize(msg)
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)

	// Cmd messages
	case codexmenu.Msg:
		return m.handleCodexAction(msg)
	case dicemenu.Msg:
		return m.handleDiceRolled(msg)
	case hintmenu.Msg:
		return m.handleHintRolled(msg)
	case oraclemenu.Msg:
		return m.handleOracleRolled(msg)

	case cmdruntui.SaveMsg:
		return m.handleSaveInput(msg)
	case cmdruntui.ScrollMsg:
		return m.handleViewportScroll(msg)
	case cmdruntui.Refresh:
		return m.handleViewRefresh(msg)
	case cmdruntui.ErrorMsg:
		return m.handleError(msg)

	case cursor.BlinkMsg:
		return m.handleCursorBlink(msg)
	}
	return m, nil
}

// HandlesEscape lets the global router forward Escape while a Codex page is open.
// func (m model) HandlesEscape() bool {
// 	// return m.codexMenu.HandlesEscape()
// 	return true
// }
