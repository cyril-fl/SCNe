package onboardforgecharacter

import (
	"fmt"
	"solopg/app/shared/services/i19n"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m model) View() tea.View {
	// content := strings.Builder{}

	// fmt.Fprintf(&content, "%s: ", transform.Capitalize(i19n.Localize("stats")))

	// for _, modifier := range m.reroll.Value {
	// 	label := i19n.Localize("stat." + string(modifier.Stat))
	// 	fmt.Fprintf(&content, "%s: %d ", label, modifier.Value)
	// }

	// content.WriteString(" | ")
	// fmt.Fprintf(&content, "%s: %d ", transform.Capitalize(i19n.Localize("attempt")), m.reroll.Attempt)

	// content.WriteString("\n\n")
	// content.WriteString(m.reroll.Options.View())

	// return tea.NewView(content.String())

	// page := page.NewPage(page.Template{
	// 	Title:    "onboarding:create_character",
	// 	Subtitle: "stats",
	// 	Body:     m.makeBody(),
	// })

	m.page.SetBody(m.makeBody())

	return m.page.GetView()
}

func (m model) makeBody() string {
	content := strings.Builder{}

	buildsChoiceStr := make([]string, len(m.reroll.Value))

	for i, modifier := range m.reroll.Value {
		key := "stat." + string(modifier.Stat)
		label := i19n.Localize(key)
		buildsChoiceStr[i] = fmt.Sprintf("%s: %d", label, modifier.Value)
	}

	content.WriteString(strings.Join(buildsChoiceStr, ", "))
	content.WriteString(strings.Join([]string{" ", i19n.Localize("attempt"), ": "}, ""))
	content.WriteString(fmt.Sprintf("%d", m.reroll.Attempt))

	content.WriteString("\n\n")
	content.WriteString(m.reroll.Options.View())

	return content.String()
}

func (m model) GetSize() *tea.WindowSizeMsg {
	return m.size
}

func (m model) SetSize(size *tea.WindowSizeMsg) {
	m.size = size
}
