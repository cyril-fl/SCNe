package combatboard

import (
	"fmt"
	"strings"

	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/gameplay/combat"
	"scne/app/cmdrun/domain/gameplay/oracle"
	"scne/app/cmdrun/tui/shortcut"
	"scne/app/shared/components/chat"
	"scne/app/shared/components/form"
	"scne/app/shared/components/models"
	"scne/app/shared/messages"
	"scne/app/shared/services/factory/millroll"
	"scne/app/shared/services/i19n"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

// - Handlers - //
// Input

func (m *model) handleKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	key := plumtypes.Key(msg.String())

	if key == plumtypes.KEY_ESC {
		return m.handleKeyEsc()
	}
	if stat, ok := getStatShortcut(string(key)); ok {
		return m.handleRoll(stat)
	}
	return m.chat.Update(msg)
}

// handleKeyEsc merges the combat thread into the main chat, then leaves the combat.
func (m *model) handleKeyEsc() tea.Cmd {
	for _, line := range m.buffer {
		m.main.Amend(line)
	}

	m.buffer = nil
	m.chat = m.makeChat()
	m.session.EndCombat()

	return form.Send[form.Close]()
}

// handleRoll rolls the selected oracle against the selected opponent, noting stat.
func (m *model) handleRoll(stat stats.Stat) tea.Cmd {
	target := m.getTarget()
	selected, ok := m.oracles.GetList().SelectedItem().(models.Item[oracle.Oracle])
	if target == nil || !ok {
		return nil
	}

	rolled := func() tea.Msg { return messages.Rolled{} }
	roll := millroll.Execute(millroll.Fight(target, stat, selected.Value()))

	return tea.Sequence(rolled, roll)
}

// Helpers

func (m *model) getTarget() *characters.Character {
	if item, ok := m.opponents.GetList().SelectedItem().(models.Item[*characters.Character]); ok {
		return item.Value()
	}
	return nil
}

// makeChat returns a chat keeping its messages in the buffer until they are merged.
func (m *model) makeChat() *chat.Chat {
	author := "System"
	if m.session.State.Player != nil {
		author = m.session.State.Player.GetName()
	}

	return chat.New(chat.Template{
		Author:  author,
		History: []string{},
		Register: func(_, message string) {
			m.buffer = append(m.buffer, message)
		},
	})
}

// refresh lists the session opponents in the side panel and greets them in the combat chat.
func (m *model) refresh() {
	l := m.opponents.GetList()
	items := make([]list.Item, 0, len(m.session.State.Opponents))
	names := make([]string, 0, len(m.session.State.Opponents))

	for _, opponent := range m.session.State.Opponents {
		name := i19n.Localize(opponent.GetName())
		items = append(items, models.NewItem(name, "", opponent))
		names = append(names, name)
	}

	l.SetItems(items)
	l.SetHeight(models.ListContentHeight(l))

	if len(m.buffer) == 0 && len(names) > 0 {
		m.chat.Amend(i19n.Localize("combat.header", map[string]any{"Opponents": strings.Join(names, ", ")}))
	}
}

// makeStatFooter returns the plain roll and one roll shortcut per stat, shown while the oracle section is active.
func makeStatFooter() []string {
	footer := []string{shortcut.THROW.String()}
	for index, stat := range stats.List() {
		footer = append(footer, shortcut.NewShortcut(shortcut.Template{
			Command: shortcut.Command(fmt.Sprintf("%s%d", shortcut.Alt, index+1)),
			Label:   "stat." + stat.String(),
		}).String())
	}
	return footer
}

// getTargetInfo writes the selected opponent name and its stats, one per line like the player stats.
func (m *model) getTargetInfo() string {
	target := m.getTarget()
	if target == nil {
		return ""
	}

	var content strings.Builder
	content.WriteString(i19n.Localize(target.GetName()))
	content.WriteString("\n")

	for _, stat := range stats.List() {
		label := i19n.Localize("stat." + stat.String())
		fmt.Fprintf(&content, "%-10s %d\n", label, target.Stats[stat])
	}

	return strings.TrimRight(content.String(), "\n")
}

// getStatShortcut returns the stat bound to key: "alt+N", or the Nth configured shortcut, rolls the Nth stat of stats.yaml.
func getStatShortcut(key string) (stats.Stat, bool) {
	shortcuts := combat.GetShortcuts()

	for index, stat := range stats.List() {
		isAlt := key == fmt.Sprintf("alt+%d", index+1)
		isAlias := index < len(shortcuts) && key == shortcuts[index]
		if isAlt || isAlias {
			return stat, true
		}
	}
	return "", false
}
