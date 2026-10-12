package codexmenu

import (
	"fmt"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/utils/transform"
	"strings"

	"charm.land/lipgloss/v2"
)

// - View - //
func (m *codexMenu) GetMenuView() string {
	title := transform.Uppercase(i19n.Localize(m.ID()))
	title = lipgloss.NewStyle().Bold(true).Render(title)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		m.GetList().View(),
	)
}

func (m *codexMenu) GetView() string {
	page := m.getCurrentPage()
	if page == nil {
		return logs.Warning("error.not_found", map[string]any{
			"Subject": transform.Capitalize(i19n.Localize("page")),
		})
	}

	return page.GetItemView()
}

func (m *codexMenuItem) GetItemView() string {
	page := m.GetPageHeader()

	if m.showForm {
		return page + m.form.View().Content
	} else {
		return page + strings.Join(m.table.Summaries(), "\n\n")
	}
}

func (m *codexMenuItem) GetPageHeader() string {
	return fmt.Sprintf("%s\n\n", i19n.Localize("sidemenu.title:page", map[string]any{
		"Submenu": transform.Capitalize(i19n.Localize(m.id)),
	}))
}

// Footer
func (m *codexMenu) GetFooter() []string {
	footer := []string{}

	if m.IsOpen() {
		footer = append(footer, i19n.Localize("cmd.shift+enter:add"))
	} else {
		footer = append(footer, i19n.Localize("cmd.shift+enter:open"))
	}

	return footer
}

// - Handlers - //
// func (m *codexMenu) HandleWindowResize(msg tea.WindowSizeMsg) tea.Cmd {
// 	if form := m.GetFormFromCurrentPage(); form != nil {
// 		form.Update(msg)
// 	}

// 	return nil
// }
