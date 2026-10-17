package gameboard

import (
	"solopg/app/cmdrun/domain/gameplay/dice"
	"solopg/app/cmdrun/domain/gameplay/hint"
	"solopg/app/cmdrun/domain/gameplay/oracle"
	"solopg/app/cmdrun/services/game"
	"solopg/app/cmdrun/tui/view/codexmenu"
	"solopg/app/cmdrun/tui/view/metadatamenu"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/navigation"
	"solopg/app/shared/messages"
	"solopg/app/shared/services/factory/millroll"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/statemachine"
	interfass "solopg/app/shared/types/interface"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type state = string

const CHAT_SECTION state = "chat"

func (m model) View() tea.View {
	msg := "404 - Not Found"
	return tea.NewView(msg)
}

// Getters & Setters
func (m *model) GetSidepanel() []string {
	return []string{
		m.menu.GetView(),
	}
}

func (m *model) SetPanelSize(size size.Size) {
	m.size = size
	m.menu.SetSize(size)
}

func (m *model) GetSection() []string {
	return []string{
		m.state.GetView(),
	}
}

func (m *model) SetSectionSize(size size.Size) {
	m.state.SetSize(size)
}

func (m model) GetFooter() []string {
	footer := []string{i19n.Localize("cmd.ctrl+s:save")}
	// TODO uniformier les comportement avec seulement alt + , ou shit + pour le side menu, dice, oracle ect on les deuw.
	footer = append(footer, m.state.GetFooter()...)
	footer = append(footer, m.menu.GetFooter()...)

	return footer
}

// Methods
func makeViewState(params UiParams) statemachine.ViewState[state] {
	params.Engine.State.Codex.EnsureInitialized()

	codexpage := codexmenu.NewFormPages(params.Engine.State.Codex)

	views := map[state]interfass.Viewable{CHAT_SECTION: makeChat(params.Engine)}

	transitions := map[state][]state{CHAT_SECTION: {}}

	for id, page := range codexpage {
		step := state(id)
		views[step] = page
		transitions[CHAT_SECTION] = append(transitions[CHAT_SECTION], step)
		transitions[step] = []state{CHAT_SECTION}
	}

	return statemachine.NewView[state](statemachine.TemplateViewState[state]{
		Base:       CHAT_SECTION,
		Transition: transitions,
		Views:      views,
	})
}

// Helpers
func makeStatView(engine *game.Engine) string {
	var content strings.Builder
	metadatamenu.GetCharacterInfo(&content, engine)
	metadatamenu.GetLocationInfo(&content, engine)
	metadatamenu.GetStatInfo(&content, engine)

	return strings.TrimRight(content.String(), "\n")
}

func makePanelPart(engine *game.Engine) []navigation.SectionInterface {
	codex := engine.State.Codex
	codex.EnsureInitialized()

	list := []navigation.SectionInterface{}
	row := []navigation.SectionInterface{
		navigation.NewSection(navigation.TemplateSection[oracle.Oracle]{
			ID:      "oracles",
			Options: oracle.ListVisible(),
			Footer:  []string{i19n.Localize("cmd.shift+enter:roll")},
			HandleEnterInput: func(value oracle.Oracle) tea.Cmd {
				return millroll.Execute(millroll.Oracles(value))
			},
			Msg: messages.Rolled{},
		}),
		navigation.NewSection(navigation.TemplateSection[dice.Dice]{
			ID:      "dice",
			Options: dice.List(),
			Footer:  []string{i19n.Localize("cmd.shift+enter:roll")},
			HandleEnterInput: func(value dice.Dice) tea.Cmd {
				return millroll.Execute(millroll.Dice(value))
			},
			Msg: messages.Rolled{},
		}),
		navigation.NewSection(navigation.TemplateSection[hint.Hint]{
			ID:      "hint",
			Options: hint.List(),
			Footer:  []string{i19n.Localize("cmd.shift+enter:roll")},
			HandleEnterInput: func(value hint.Hint) tea.Cmd {
				return millroll.Execute(millroll.Hint(value))
			},
			Msg: messages.Rolled{},
		}),
		codexmenu.NewSection(false),
	}

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
