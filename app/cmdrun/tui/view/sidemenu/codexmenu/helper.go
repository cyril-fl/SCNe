package codexmenu

import (
	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/cmdrun/domain/card/characters/classes"
	"solopg/app/cmdrun/domain/card/characters/races"
	"solopg/app/cmdrun/domain/gameplay/oracle"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/form/field"
	"solopg/app/shared/components/models"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/config"

	tea "charm.land/bubbletea/v2"
)

// - Side panel - //
func makeCodexList(params CodexMenuParams) []*codexMenuItem {
	pages := []*codexMenuItem{
		newMenuItem(sideMenuItemsParams{
			id:    "codex.locations",
			table: params.Codex.LocationsTable,
			form:  getLocationForm,
		}),
		newMenuItem(sideMenuItemsParams{
			id:    "codex.npcs",
			table: params.Codex.NpcsTable,
			form:  getNpcForm,
		}),
		newMenuItem(sideMenuItemsParams{
			id:    "codex.beastiary",
			table: params.Codex.BeastiaryTable,
			form:  getBeastiaryForm,
		}),
		newMenuItem(sideMenuItemsParams{
			id:    "codex.objects",
			table: params.Codex.ObjectsTable,
			form:  getObjectForm,
		}),
		newMenuItem(sideMenuItemsParams{
			id:    "codex.objectives",
			table: params.Codex.ObjectifsTable,
			form:  getObjectifForm,
		}),
	}

	return pages
}

// - Getters & Setters - //
// Codex menu
func (m *codexMenu) getCurrentPage() *codexMenuItem {
	for _, item := range m.GetList().Items() {

		if page, ok := item.(models.Item[*codexMenuItem]); !ok {
			continue
		} else if page.Value().isOpen {
			return page.Value()
		}
	}
	return nil
}

func (m *codexMenu) getSelectedPage() *codexMenuItem {
	if item, ok := m.GetList().SelectedItem().(models.Item[*codexMenuItem]); ok {
		return item.Value()
	}
	return nil
}

// Codex menu iteme
func (m *codexMenuItem) setOpen(open bool) {
	if open {
		m.isOpen = true
		return
	}

	m.isOpen = false
	m.showForm = false
}

func (m *codexMenuItem) setShowForm(show bool) {
	state := "closed"
	if show {
		state = "opened"
	}

	logs.SilentInfo("navigation."+state, map[string]any{
		"Subject": "form",
		"Value":   m.id,
	})

	m.showForm = show
}

func (m *codexMenuItem) toggleShowForm() {
	m.showForm = !m.showForm
}

// - Helper - //
// Input
func (m *codexMenu) handleKeyEsc(msg tea.KeyPressMsg) tea.Cmd {
	currentPage := m.getCurrentPage()

	if currentPage == nil {
		return nil
	}

	if currentPage.showForm {
		currentPage.toggleShowForm()
		return sendMsg()
	}

	if currentPage.isOpen {
		currentPage.isOpen = false
		return sendMsg()
	}

	return sendMsg()
}

func (m *codexMenu) handleKeyShiftEnter(msg tea.KeyPressMsg) tea.Cmd {
	currentPage := m.getCurrentPage()

	if currentPage != nil && m.IsOpen() {
		if currentPage.showForm {
			return form.SendMsg[form.Validate]()
		}

		currentPage.form = currentPage.newForm()
		currentPage.setShowForm(true)
		return sendMsg()
	}

	if selectedPage := m.getSelectedPage(); selectedPage != nil {
		m.handleOpenPage(selectedPage)
	}

	return sendMsg()
}

// Page
func (m *codexMenu) handleOpenPage(selected *codexMenuItem) {
	for _, item := range m.GetList().Items() {
		item, ok := item.(models.Item[*codexMenuItem])

		if page := item.Value(); ok && page != selected {

			logs.SilentInfo("navigation.closed", map[string]any{
				"Subject": "page",
				"Value":   page.id,
			})

			page.setOpen(false)
		}
	}

	logs.SilentInfo("navigation.open", map[string]any{
		"Subject": "page",
		"Value":   selected.id,
	})

	selected.setOpen(true)
}

func (m *codexMenu) updatePageOnRedirection() {
	currentPage := m.getCurrentPage()

	if currentPage != nil && m.IsOpen() {
		selected := m.getSelectedPage()
		m.handleOpenPage(selected)
	}
}

// Form
func getFormRaceOptions() []models.Item[string] {
	var raceItems []models.Item[string]
	for _, i := range races.List() {
		raceItems = append(raceItems, models.NewItem(i19n.Localize(i.GetName()), "", i.GetName()))
	}
	return raceItems
}

func getFormClassOptions() []models.Item[string] {
	var classItems []models.Item[string]
	for _, i := range classes.List() {
		classItems = append(classItems, models.NewItem(i19n.Localize(i.GetName()), "", i.GetName()))
	}
	classItems = append(classItems, models.NewItem("None", "", ""))
	return classItems
}

func getFormOracle() []models.Item[string] {
	list, err := oracle.ListFromFolder(config.Current.Documents.Folders.Encounters)
	if err != nil {
		return nil
	}

	var oracleItems []models.Item[string]
	for _, i := range list {
		oracleItems = append(oracleItems, models.NewItem(i19n.Localize("oracles."+i), "", i))
	}

	return oracleItems
}

func getFormStatsFields() []form.Field {
	fields := []form.Field{}

	for _, i := range stats.List() {
		f := field.TextField(field.TextTemplate[string]{
			ID:           i.String(),
			Label:        i19n.Localize(i.String()),
			Validator:    nil,
			Required:     false,
			Defaultvalue: "",
		})
		fields = append(fields, f)
	}

	return fields
}

func (m *codexMenu) handleFormSubmit(msg tea.Msg) tea.Cmd {
	currentPage := m.getCurrentPage()
	if currentPage == nil {
		return nil
	}

	currentForm := currentPage.form
	if currentForm == nil {
		return nil
	}

	currentForm.Submit()

	if currentForm.HasErr() {
		return form.SendErrorMsg(currentForm.GetErr())
	}

	logs.SilentInfo("form.submitted", map[string]any{
		"Form": currentPage.id,
	})

	return m.handleFormPost()
}

func (m *codexMenu) handleFormPost() tea.Cmd {
	currentPage := m.getCurrentPage()
	if currentPage == nil {
		return nil
	}

	currentForm := currentPage.form
	if currentForm == nil {
		return nil
	}

	err := currentPage.table.AddFromMappedValues(currentForm.GetValuesAsMappedString())
	if err != nil {
		currentPage.form.SetErr(err)
		return form.SendErrorMsg(err)
	}

	currentPage.setShowForm(false)

	logs.SilentSuccess("form.post", map[string]any{
		"Form": currentPage.id,
	})

	return sendMsg()
}

// func (m *codexMenu) handleFormScroll(msg tea.KeyPressMsg) tea.Cmd {
// 	form := m.GetFormFromCurrentPage()
// 	if form == nil {
// 		return  nil
// 	}

// 	currentPage := m.getCurrentPage()
// 	pageHeader := lipgloss.NewStyle().
// 		// Width(params.Viewport.Width()).
// 		Render(currentPage.GetPageHeader())

// 	// viewTop, viewHight := form.GetFieldScrollArea(params.Viewport)
// 	// viewHight += lipgloss.Height(pageHeader) ;

// 	return  sharedtui.SendScrollMsg(viewTop, viewHight)
// }

func (m *codexMenu) handleFormError(msg tea.Msg) tea.Cmd {
	return sendMsg()
}

func (m *codexMenu) delegateInputToForm(msg tea.Msg) tea.Cmd {
	currentPage := m.getCurrentPage()

	if currentPage == nil {
		return nil
	}

	if !currentPage.showForm {
		return nil
	}

	_, cmd := currentPage.form.Update(msg)

	return tea.Batch(
		cmd,
		sendMsg(),
	)
}
