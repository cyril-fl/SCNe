package onboardlocation

import (
	"fmt"
	"solopg/app/shared/services/i19n"
	sharedtui "solopg/app/shared/tui"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// TODO Ui ici bug !
func (m model) View() tea.View {
	// content := strings.Builder{}

	// fmt.Fprintf(&content, "%s: %s", transform.Capitalize(i19n.Localize("location")), i19n.Localize(m.reroll.Value.Name))
	// content.WriteString(" | ")
	// fmt.Fprintf(&content, "%s: %d", transform.Capitalize(i19n.Localize("attempt")), m.reroll.Attempt)

	// if m.reroll.Value.Description.Description != sharedtui.KEY_EMPTY {
	// 	content.WriteString("\n")
	// 	content.WriteString(i19n.Localize(m.reroll.Value.Description.Description))
	// }

	// content.WriteString("\n\n")
	// content.WriteString(m.reroll.Options.View())

	// return tea.NewView(content.String())

	// TODO a un moment autour d'ici il y a yb proble avec la view et un "scrol" se fait
	// le probleme avec dans un rerol, deha ic avec le build

	// page := page.NewPage(page.Template{
	// 	Title:    "onboarding.character:create",
	// 	Subtitle: "location",
	// 	Body:     m.makeBody(),
	// })

	m.page.SetBody(m.makeBody())

	return m.page.GetView()
}

func (m model) makeBody() string {
	content := strings.Builder{}

	content.WriteString(m.reroll.Value.Name)
	content.WriteString(strings.Join([]string{" ", i19n.Localize("attempt"), ": "}, ""))
	content.WriteString(fmt.Sprintf("%d", m.reroll.Attempt))

	if m.reroll.Value.Description.Description != sharedtui.KEY_EMPTY {
		content.WriteString("\n")
		content.WriteString(m.reroll.Value.Description.Description)
	}

	content.WriteString("\n\n")
	content.WriteString(m.reroll.Options.View())
	return content.String()
}

func (m model) HandlesEscape() bool {
	// return m.codexMenu.HandlesEscape()
	return true
}
