package combatboard

import (
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/gameplay/oracle"
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/shared/components/chat"
	"scne/app/shared/components/navigation"
	"scne/app/shared/messages"
	"scne/app/shared/services/factory/millroll"

	tea "charm.land/bubbletea/v2"
)

// - Model - //
type model struct {
	session   *session.GameSession
	main      *chat.Chat
	chat      *chat.Chat
	opponents *navigation.FocusableSection
	oracles   *navigation.FocusableSection
	target    *navigation.Section
	buffer    []string
}

type UiParams struct {
	Session *session.GameSession
	// Main is the gameboard chat the combat thread is merged into on leave.
	Main *chat.Chat
}

func NewModel(params UiParams) *model {
	m := &model{
		session: params.Session,
		main:    params.Main,
		opponents: navigation.NewFocusableSection(navigation.TemplateFocusableSection[*characters.Character]{
			ID: "opponents",
		}),
		oracles: navigation.NewFocusableSection(navigation.TemplateFocusableSection[oracle.Oracle]{
			ID:      "oracles",
			Options: oracle.ListVisible(),
			Footer:  makeStatFooter(),
			HandleEnterInput: func(value oracle.Oracle) tea.Cmd {
				return millroll.Execute(millroll.Oracles(value))
			},
			Msg: messages.Rolled{},
		}),
	}

	m.chat = m.makeChat()
	m.target = navigation.NewSection(navigation.TemplateSection{
		ID:      "target",
		Content: m.getTargetInfo,
	})

	return m
}

// Getters & Setters

// GetOpponentSection returns the side panel section listing the opponents; its selection is the target.
func (m *model) GetOpponentSection() *navigation.FocusableSection {
	return m.opponents
}

// GetTargetSection returns the side panel block showing the selected opponent stats.
func (m *model) GetTargetSection() *navigation.Section {
	return m.target
}

// GetOracleSection returns the side panel oracle section; its selection is rolled by the stat shortcuts shown in its footer.
func (m *model) GetOracleSection() *navigation.FocusableSection {
	return m.oracles
}

// Methods

// Update rolls the selected oracle on alt+N, merges and leaves on esc, and forwards other messages to the combat chat.
func (m *model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	}
	return m.chat.Update(msg)
}

// Amend adds content to the combat chat.
func (m *model) Amend(content string) {
	m.chat.Amend(content)
}
