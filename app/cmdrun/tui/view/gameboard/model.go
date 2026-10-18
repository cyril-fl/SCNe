package gameboard

import (
	"solopg/app/cmdrun/domain/gameplay/session"
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
)

// - Model - //
type model struct {
	plumprimitive.Fallible

	state   statemachine.ViewState[state]
	menu    *navigation.Navigation
	session *session.GameSession
	save    func() error
	size    size.Size
}

type UiParams struct {
	Session *session.GameSession
	OnSave  func() error
}

func NewModel(params UiParams) *model {
	return &model{
		state: makeViewState(params),
		menu: navigation.New(navigation.Template{
			Static:     []string{makeStatView(params.Session)},
			Scrollable: makePanelPart(params.Session),
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

	var menucmd tea.Cmd
	if isChatSection || isCloseForm {
		menucmd = m.menu.HandleUpdate(msg)
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
