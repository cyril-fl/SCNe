package gameboard

import (
	"solopg/app/cmdrun/services/game"
	cmdruntui "solopg/app/cmdrun/tui"
	"solopg/app/cmdrun/tui/view/sidemenu"
	"solopg/app/cmdrun/tui/view/sidemenu/codexmenu"
	"solopg/app/cmdrun/tui/view/sidemenu/dicemenu"
	"solopg/app/cmdrun/tui/view/sidemenu/hintmenu"
	"solopg/app/cmdrun/tui/view/sidemenu/oraclemenu"
	"solopg/app/cmdrun/types/direction"
	"solopg/app/cmdrun/types/size"
	"solopg/app/cmdrun/types/viewoptions"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"
	"strings"
	"time"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// - Getters & Setters - //
func (m *model) getMenuActiveElement() sidemenu.MenuItem {
	return m.menu[m.activeMenuIndex]
}

func (m *model) isLastMenuElement() bool {
	return m.activeMenuIndex == len(m.menu)-1
}

func (m *model) isFirstMenuElement() bool {
	return m.activeMenuIndex == 0
}

func (m *model) isMenuActiveElementOpen() bool {
	return m.getMenuActiveElement().IsOpen()
}

// - Handlers - //
// Input
func (m *model) handleEnterInput() {
	if m.isMenuActiveElementOpen() {
		return
	}
	input := m.textarea.Value()
	if input == "" {
		return
	}

	amendChat(m, m.author, input)

	m.textarea.Reset()
	refreshViewport(m, viewoptions.RefreshOption{ScrollBottom: true})
}

func (m *model) handleSaveInput(msg cmdruntui.SaveMsg) (*model, tea.Cmd) {
	m.SetErr(msg.Err)

	if m.HasErr() {
		err := logs.Error("error.unexpected:save", map[string]any{
			"Error": msg.Err,
		})

		m.chat = append(m.chat, err.Error())

	} else {
		log := logs.Success("campaign:success", map[string]any{
			"Time": time.Now().Format("2006-01-02 15:04:05")})

		m.chat = append(m.chat, log)
	}

	return refreshViewport(m, viewoptions.RefreshOption{ScrollBottom: true})
}

func handleDefaultInput(m *model, msg tea.Msg) (*model, tea.Cmd) {
	var cmd tea.Cmd

	if m.isMenuActiveElementOpen() {
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m *model) handleCursorBlink(msg cursor.BlinkMsg) (*model, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *model) handleKeyPress(msg tea.KeyPressMsg) (*model, tea.Cmd) {
	switch msg.String() {
	case sharedtui.KEY_ENTER:
		m.handleEnterInput()
	case sharedtui.KEY_UP, sharedtui.KEY_DOWN:
		m.handleDirectionInput(msg)
	default:
		handleDefaultInput(m, msg)
	}
	return m.handleCommand(msg)
}

func (m *model) handleCommand(msg tea.KeyPressMsg) (*model, tea.Cmd) {
	switch msg.String() {
	case sharedtui.CMD_CTRL_S:
		return m, cmdruntui.SendSaveMsg(m.save)
	}

	return refreshViewport(m, viewoptions.RefreshOption{ScrollBottom: true})
}

func (m *model) handleError(msg cmdruntui.ErrorMsg) (*model, tea.Cmd) {
	m.SetErr(msg.Err)
	return refreshViewport(m, viewoptions.RefreshOption{ScrollBottom: true})
}

// - Side Menu - //
// Codex Action
func (m *model) handleCodexAction(msg codexmenu.Msg) (*model, tea.Cmd) {
	_ = msg

	return refreshViewport(m, viewoptions.RefreshOption{Positionreset: false})
}

// Dice Rolled
func (m *model) handleDiceRolled(msg dicemenu.Msg) (*model, tea.Cmd) {
	message := i19n.Localize("dice.roll:result", map[string]any{
		"Dice":  msg.Dice,
		"Value": msg.Value,
	})

	// -- Advent journal
	amendChat(m, "Dice", message)

	return refreshViewport(m, viewoptions.RefreshOption{ScrollBottom: true})
}

// Hint Rolled
func (m *model) handleHintRolled(msg hintmenu.Msg) (*model, tea.Cmd) {
	localized := make([]string, 0, len(msg.Result))
	for _, hint := range msg.Result {
		localized = append(localized, i19n.Localize(hint))
	}

	message := i19n.Localize("hint.roll:result", map[string]any{
		"Value": strings.Join(localized, ", "),
	})

	amendChat(m, "Hint", message)

	return refreshViewport(m, viewoptions.RefreshOption{ScrollBottom: true})
}

// Oracle Rolled
func (m *model) handleOracleRolled(msg oraclemenu.Msg) (*model, tea.Cmd) {
	message := i19n.Localize("oracle.roll:result", map[string]any{
		"Roll":     msg.Result.Roll,
		"Value":    msg.Result.Result,
		"Critical": msg.Result.Critical,
	})

	amendChat(m, "Oracle", message)

	return refreshViewport(m, viewoptions.RefreshOption{ScrollBottom: true})
}

// Side Menu Direction
func (m *model) handleDirectionInput(key tea.KeyPressMsg) {
	currentMenu := m.getMenuActiveElement()

	direction := sidemenu.HandleKeyArrow(sidemenu.Context{
		Menu:             currentMenu,
		CurrentMenuIndex: m.activeMenuIndex,
		SibblingCount:    len(m.menu),
	}, key)

	m.updateMenuDirection(direction)

	currentMenu.HandleDirectionInput(direction)
}

func (m *model) updateMenuDirection(msg direction.Direction) {
	previousMenuIndex := m.activeMenuIndex
	menuLength := len(m.menu)

	switch msg {
	case direction.PREVIOUS:
		if m.activeMenuIndex > 0 {
			m.activeMenuIndex--
		}
	case direction.NEXT:
		if m.activeMenuIndex < menuLength-1 {
			m.activeMenuIndex++
		}
	}

	if previousMenuIndex != m.activeMenuIndex {
		m.menu[previousMenuIndex].SetFocus(false)
		m.menu[m.activeMenuIndex].SetFocus(true)
	}
}

// - Window - //
func (m *model) handleWindowResize(msg tea.WindowSizeMsg) (*model, tea.Cmd) {
	refreshLayout(m, msg)
	refreshMenuElement(m, msg)

	return m, nil
}

func (m *model) handleViewportScroll(msg tea.Msg) (*model, tea.Cmd) {
	scrollMsg, ok := msg.(cmdruntui.ScrollMsg)
	if !ok {
		return m, nil
	}

	top := scrollMsg.GetTop()
	hight := scrollMsg.GetHeight()
	bottom := scrollMsg.GetBottom()

	v := &m.viewport
	vHeight := v.Height()
	vOffset := v.YOffset()

	if rendered := vHeight > 0; !rendered {
		return m, nil
	} else if scrollUp := top < vOffset; scrollUp {
		v.SetYOffset(top)
	} else if scrollDown := bottom > vOffset+vHeight; scrollDown {
		v.SetYOffset(top + hight - vHeight)
	}

	return refreshViewport(m, viewoptions.RefreshOption{Positionreset: false})
}

func (m model) handleViewRefresh(msg cmdruntui.Refresh) (tea.Model, tea.Cmd) {
	return refreshViewport(&m, msg.Option)
}

// - Helper - //
// Init
func initTextarea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = i19n.Localize("chat.input:placeholder")
	ta.SetVirtualCursor(false)
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(s)

	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return ta
}

func initViewport() viewport.Model {
	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	return vp
}

func initJournal(engine *game.Engine) []string {
	raw := make([]string, 0, len(engine.State.AdventureLog.Entries))

	for _, entry := range engine.State.AdventureLog.Entries {
		raw = append(raw, entry.String())
	}

	return raw
}

func initSideMenu(engine *game.Engine) []sidemenu.MenuItem {
	size := size.New(size.Template{
		Width:  PANEL_WHIDTH - 4,
		Height: ORACLE_HEIGHT + 5,
	})

	codex := engine.State.Codex
	codex.EnsureInitialized()

	return []sidemenu.MenuItem{
		oraclemenu.NewSideMenu(size, true),
		dicemenu.NewSideMenu(size, false),
		hintmenu.NewSideMenu(size, false),
		codexmenu.NewSideMenu(codexmenu.CodexMenuParams{
			Size:  size,
			Codex: codex,
		}, false),
	}
}

func initAuthor(engine *game.Engine) string {
	author := "System"
	if engine.State.Player != nil {
		author = engine.State.Player.Name
	}
	return author
}

// Refresh
func refreshMenuElement(m *model, msg tea.WindowSizeMsg) {
	activeEl := m.getMenuActiveElement()
	activeEl.HandleWindowResize(msg)
}

func amendChat(m *model, author, message string) {
	adventureJournal := m.engine.GetAdventureJournal()
	adventureJournal.AddEntry(author, message)

	authorstyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("5")).
		Render(m.author + ": ")

	cloredmsg := authorstyle + message
	m.chat = append(m.chat, cloredmsg)
}
