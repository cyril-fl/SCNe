package gameboard

import (
	"scne/app/cmdrun/domain/gameplay/dice"
	"scne/app/cmdrun/domain/gameplay/hint"
	"scne/app/cmdrun/domain/gameplay/oracle"
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/cmdrun/tui/shortcut"
	"scne/app/cmdrun/tui/view/menuaction"
	"scne/app/cmdrun/tui/view/menucodex"
	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/navigation"
	"scne/app/shared/messages"
	"scne/app/shared/services/factory/millroll"
	"scne/app/shared/services/statemachine"
	interfass "scne/app/shared/types/interface"

	tea "charm.land/bubbletea/v2"
)

type state = string

const CHAT_SECTION state = "chat"

const COMBAT_SECTION state = menuaction.COMBAT_PAGE

func (m model) View() tea.View {
	msg := "404 - Not Found"
	return tea.NewView(msg)
}

// Getters & Setters
func (m *model) GetSidepanel() []string {
	return []string{
		m.activeMenu().GetView(),
	}
}

func (m *model) SetPanelSize(size size.Size) {
	m.size = size
	m.menu.SetSize(size)
	m.combatMenu.SetSize(size)
}

// activeMenu returns the combat panel while in combat, the default panel otherwise.
func (m *model) activeMenu() *navigation.Navigation {
	if m.state.Current() == COMBAT_SECTION {
		return m.combatMenu
	}
	return m.menu
}

func (m *model) GetSection() []string {
	return []string{
		m.state.GetView(),
	}
}

func (m *model) SetSectionSize(size size.Size) {
	m.state.SetSize(size)
}

func (m *model) GetFooter() []string {
	footer := []string{}

	if m.state.Current() == CHAT_SECTION {
		footer = append(footer, shortcut.SAVE.String(), shortcut.NAVIGATE.String())
	}

	footer = append(footer, m.state.GetFooter()...)
	footer = append(footer, m.activeMenu().GetFooter()...)
	return footer
}

// Methods
func makeViewState(params UiParams, mainchat interfass.Viewable, combatpage interfass.Viewable) statemachine.ViewState[state] {
	params.Session.State.Codex.EnsureInitialized()

	codexpage := menucodex.NewFormPages(params.Session)

	views := map[state]interfass.Viewable{
		CHAT_SECTION:   mainchat,
		COMBAT_SECTION: combatpage,
	}

	transitions := map[state][]state{
		CHAT_SECTION:   {COMBAT_SECTION},
		COMBAT_SECTION: {CHAT_SECTION},
	}

	for id, page := range codexpage {
		step := state(id)
		views[step] = page
		transitions[CHAT_SECTION] = append(transitions[CHAT_SECTION], step)
		transitions[step] = []state{CHAT_SECTION}
	}

	actionpage := menuaction.NewFormPages(params.Session)

	for id, page := range actionpage {
		step := state(id)
		views[step] = page
		transitions[CHAT_SECTION] = append(transitions[CHAT_SECTION], step)
		transitions[step] = []state{CHAT_SECTION}
	}

	return statemachine.NewView(statemachine.TemplateViewState[state]{
		Base:       CHAT_SECTION,
		Transition: transitions,
		Views:      views,
	})
}

// Helpers
func makePanelPart(gameSession *session.GameSession) []navigation.FocusableSectionInterface {
	codex := gameSession.State.Codex
	codex.EnsureInitialized()

	row := []navigation.FocusableSectionInterface{makeOracleSection()}
	row = append(row, makeRollSections()...)
	row = append(row, menuaction.NewSection(false), menucodex.NewSection(false))

	return focusFirst(row)
}

// makeCombatPanelPart lists the opponents and the combat oracles, in place of the actions and the codex.
func makeCombatPanelPart(opponents, oracles navigation.FocusableSectionInterface) []navigation.FocusableSectionInterface {
	row := []navigation.FocusableSectionInterface{opponents, oracles}
	row = append(row, makeRollSections()...)

	return focusFirst(row)
}

// makeOracleSection returns the oracle section of the default panel; the combat panel uses its own.
func makeOracleSection() navigation.FocusableSectionInterface {
	return navigation.NewFocusableSection(navigation.TemplateFocusableSection[oracle.Oracle]{
		ID:      "oracles",
		Options: oracle.ListVisible(),
		Footer:  []string{shortcut.THROW.String()},
		HandleEnterInput: func(value oracle.Oracle) tea.Cmd {
			return millroll.Execute(millroll.Oracles(value))
		},
		Msg: messages.Rolled{},
	})
}

// makeRollSections returns the dice and hint sections, shared by both panels.
func makeRollSections() []navigation.FocusableSectionInterface {
	return []navigation.FocusableSectionInterface{
		navigation.NewFocusableSection(navigation.TemplateFocusableSection[dice.Dice]{
			ID:      "dice",
			Options: dice.List(),
			Footer:  []string{shortcut.THROW.String()},
			HandleEnterInput: func(value dice.Dice) tea.Cmd {
				return millroll.Execute(millroll.Dice(value))
			},
			Msg: messages.Rolled{},
		}),
		navigation.NewFocusableSection(navigation.TemplateFocusableSection[hint.Hint]{
			ID:      "hint",
			Options: hint.List(),
			Footer:  []string{shortcut.THROW.String()},
			HandleEnterInput: func(value hint.Hint) tea.Cmd {
				return millroll.Execute(millroll.Hint(value))
			},
			Msg: messages.Rolled{},
		}),
	}
}

func focusFirst(row []navigation.FocusableSectionInterface) []navigation.FocusableSectionInterface {
	list := []navigation.FocusableSectionInterface{}

	for i, item := range row {
		if i > 0 {
			item.SetFocus(false)
		} else {
			item.SetFocus(true)
		}

		list = append(list, item)
	}

	return list
}
