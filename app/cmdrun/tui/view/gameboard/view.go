package gameboard

import (
	"solopg/app/cmdrun/domain/gameplay/dice"
	"solopg/app/cmdrun/domain/gameplay/hint"
	"solopg/app/cmdrun/domain/gameplay/oracle"
	"solopg/app/cmdrun/services/factory/milldice"
	"solopg/app/cmdrun/services/factory/millhint"
	"solopg/app/cmdrun/services/factory/milloracle"
	"solopg/app/cmdrun/services/factory/roll"
	"solopg/app/cmdrun/services/game"
	"solopg/app/cmdrun/tui/view/sidemenu/codexmenu"
	"solopg/app/cmdrun/tui/view/sidemenu/metadatamenu"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/navigation"
	"solopg/app/shared/services/i19n"
	"strings"

	tea "charm.land/bubbletea/v2"
)

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
	return append(footer, m.state.getFooter()...)
}

// Methods

// Helpers
func makeStatView(engine *game.Engine) string {
	var content strings.Builder
	metadatamenu.GetCharacterInfo(&content, engine)
	metadatamenu.GetLocationInfo(&content, engine)
	metadatamenu.GetStatInfo(&content, engine)

	return strings.TrimRight(content.String(), "\n")
}

func makePanelPart(engine *game.Engine) []navigation.Item {
	codex := engine.State.Codex
	codex.EnsureInitialized()

	list := []navigation.Item{}
	row := []navigation.Item{
		navigation.NewMenu(navigation.TemplateMenu[oracle.Oracle]{
			ID:      "oracles",
			Options: oracle.ListVisible(),
			Footer:  []string{i19n.Localize("cmd.shift+enter:roll")},
			HandleEnterInput: func(value oracle.Oracle) tea.Cmd {
				return roll.Execute(milloracle.New(value))
			},
		}),
		navigation.NewMenu(navigation.TemplateMenu[dice.Dice]{
			ID:      "dice",
			Options: dice.List(),
			Footer:  []string{i19n.Localize("cmd.shift+enter:roll")},
			HandleEnterInput: func(value dice.Dice) tea.Cmd {
				return roll.Execute(milldice.New(value))
			},
		}),
		navigation.NewMenu(navigation.TemplateMenu[hint.Hint]{
			ID:      "hint",
			Options: hint.List(),
			Footer:  []string{i19n.Localize("cmd.shift+enter:roll")},
			HandleEnterInput: func(value hint.Hint) tea.Cmd {
				return roll.Execute(millhint.New(value))
			},
		}),
		codexmenu.NewSideMenu(codexmenu.CodexMenuParams{
			Size:  size.Nil(),
			Codex: codex,
		}, false),
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
