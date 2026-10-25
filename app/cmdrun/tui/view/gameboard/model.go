package gameboard

import (
	"solopg/app/cmdrun/domain/gameplay/session"
	"solopg/app/cmdrun/tui/view/combatboard"
	"solopg/app/cmdrun/tui/view/menumetadata"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/chat"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/navigation"
	"solopg/app/shared/messages"
	"solopg/app/shared/services/statemachine"
	"time"

	plumprimitive "github.com/cyril-fl/plum/primitive"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

// - Model - //
type model struct {
	plumprimitive.Fallible

	state      statemachine.ViewState[state]
	menu       *navigation.Navigation
	combatMenu *navigation.Navigation
	session    *session.GameSession
	save       func() error
	size       size.Size
}

type UiParams struct {
	Session *session.GameSession
	OnSave  func() error
}

func NewModel(params UiParams) *model {
	mainchat := makeChat(params.Session)
	combatpage := combatboard.NewModel(combatboard.UiParams{
		Session: params.Session,
		Main:    mainchat,
	})

	return &model{
		state: makeViewState(params, mainchat, combatpage),
		menu: navigation.New(navigation.Template{
			Static:     []navigation.SectionInterface{menumetadata.NewSection(params.Session)},
			Scrollable: makePanelPart(params.Session),
		}),
		combatMenu: navigation.New(navigation.Template{
			Static:     []navigation.SectionInterface{menumetadata.NewSection(params.Session), combatpage.GetTargetSection()},
			Scrollable: makeCombatPanelPart(combatpage.GetOpponentSection(), combatpage.GetOracleSection()),
		}),
		session: params.Session,
		save:    params.OnSave,
	}
}

// Getters & Setters

// Methods
func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	statecmd := m.state.Update(msg)

	_, ok := msg.(form.Close)
	isCloseForm := ok
	isChatSection := m.state.Current() == CHAT_SECTION
	isCombatSection := m.state.Current() == COMBAT_SECTION

	// form.Close always reaches the default panel so its sections forget the closed page.
	menu := m.activeMenu()
	if isCloseForm {
		menu = m.menu
	}

	var menucmd tea.Cmd
	if menumsg, ok := toMenuMsg(msg); ok && (isChatSection || isCombatSection || isCloseForm) {
		menucmd = menu.HandleUpdate(menumsg)
	}

	next, cmd := m.handleUpdate(msg)

	return next, tea.Sequence(menucmd, statecmd, cmd)
}

func (m *model) handleUpdate(msg tea.Msg) (*model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	case form.Open:
		m.state.HandleTransition(state(msg.ID))
	case form.Close:
		m.state.HandleTransition(CHAT_SECTION)
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

// toMenuMsg keeps plain arrows for the chat textarea and maps alt+arrows to menu navigation.
func toMenuMsg(msg tea.Msg) (tea.Msg, bool) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return msg, true
	}

	switch plumtypes.Key(key.String()) {
	case plumtypes.KEY_UP, plumtypes.KEY_DOWN:
		return nil, false
	case plumtypes.KEY_ALT_UP:
		return tea.KeyPressMsg{Code: tea.KeyUp}, true
	case plumtypes.KEY_ALT_DOWN:
		return tea.KeyPressMsg{Code: tea.KeyDown}, true
	}

	return msg, true
}

func makeChat(gameSession *session.GameSession) *chat.Chat {
	return chat.New(chat.Template{
		Author:  initAuthor(gameSession),
		History: initJournal(gameSession),
		Register: func(author, message string) {
			gameSession.GetAdventureJournal().AddEntry(author, message)
		},
	})
}

func initAuthor(gameSession *session.GameSession) string {
	author := "System"
	if gameSession.State.Player != nil {
		author = gameSession.State.Player.GetName()
	}
	return author
}

func initJournal(gameSession *session.GameSession) []string {
	raw := make([]string, 0, len(gameSession.State.AdventureLog.Entries))
	for _, entry := range gameSession.State.AdventureLog.Entries {
		raw = append(raw, entry.String())
	}
	return raw
}

// HandlesEscape lets the global router forward Escape while a Codex page is open.
func (m model) HandlesEscape() bool {
	return true
}

// Helpers
