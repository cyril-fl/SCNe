package combatboard

import (
	"strings"

	"solopg/app/cmdrun/domain/gameplay/combat"
	"solopg/app/cmdrun/tui/shortcut"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/models"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/style"

	"charm.land/lipgloss/v2"
)

// Getters & Setters

// GetView renders the fighters, the player actions (or the outcome) and the latest log lines.
func (m *model) GetView() string {
	current := m.session.State.Combat
	if current == nil {
		return ""
	}

	header := lipgloss.JoinVertical(lipgloss.Left,
		style.Title.Render(i19n.Localize("combat.round", map[string]any{"Round": current.Round + 1})),
		renderFighter(current.Opponent),
		renderFighter(current.Player),
		"",
	)

	body := style.Accent.Render(i19n.Localize(current.GetOutcome()))
	if !current.IsOver() {
		list := m.GetList()
		body = list.View()
	}

	used := lipgloss.Height(header) + lipgloss.Height(body) + 1
	log := style.Muted.Render(strings.Join(lastLines(m.log, m.size.GetHeight()-used), "\n"))

	return lipgloss.JoinVertical(lipgloss.Left, header, body, "", log)
}

// GetFooter returns play and leave while the combat runs, leave once it is over.
func (m *model) GetFooter() []string {
	current := m.session.State.Combat
	if current == nil || current.IsOver() {
		return []string{shortcut.LEAVE.String()}
	}
	return []string{shortcut.PLAY.String(), shortcut.LEAVE.String()}
}

// SetSize stores s and reloads the player actions; the gameboard calls it when the page is entered.
func (m *model) SetSize(s size.Size) {
	m.size = s
	m.refresh()

	list := m.GetList()
	m.SetListSize(size.Template{
		Width:  s.GetWidth(),
		Height: models.ListContentHeight(list),
	})
}

// Helpers

func renderFighter(fighter *combat.Fighter) string {
	return i19n.Localize("combat.fighter", map[string]any{
		"Fighter": fighter.Character.GetName(),
		"Health":  fighter.GetHealth(),
	})
}

// lastLines returns the last count lines, at least one.
func lastLines(lines []string, count int) []string {
	count = max(1, count)
	if len(lines) <= count {
		return lines
	}
	return lines[len(lines)-count:]
}
