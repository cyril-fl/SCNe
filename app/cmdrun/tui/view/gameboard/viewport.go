package gameboard

// - Viewport - //

// - Helper - //
// Refresh
// func refreshViewport(m *model, params viewoptions.RefreshOption) (*model, tea.Cmd) {
// 	// content := m.getContent()

// 	// m.viewport.SetContent(content)

// 	if params.Positionreset {
// 		// m.viewport.GotoTop()
// 	} else if params.ScrollBottom {
// 		// m.viewport.GotoBottom()
// 	}

// 	return m, nil
// }

// func refreshLayout(m *model, msg tea.WindowSizeMsg) (*model, tea.Cmd) {
// 	// chatWidth := max(0, msg.Width-PANEL_WHIDTH-PANEL_GAP)

// 	// m.viewport.SetHeight(max(0, msg.Height-m.textarea.Height()-1))
// 	// m.viewport.SetWidth(chatWidth)
// 	// m.textarea.SetWidth(chatWidth)

// 	// if !m.isMenuActiveElementOpen() {
// 	// 	// m.viewport.SetContent(m.getContent())
// 	// }

// 	// m.viewport.GotoBottom()

// 	return m, nil
// }

// Render
// func getCursor(m model) *tea.Cursor {
// 	// cursor := m.textarea.Cursor()
// 	// cursor.Y += lipgloss.Height(m.viewport.View())

// 	if item := m.getMenuActiveElement(); item != nil && item.IsOpen() {
// 		// cursor = nil
// 	}

// 	return cursor
// }
