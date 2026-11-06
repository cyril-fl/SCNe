package combatboard

import (
	"scne/app/cmdrun/types/size"

	tea "charm.land/bubbletea/v2"
)

// plugin adds a shortcut to the combat; a plugin with a view takes over the page.
type plugin interface {
	HandleUpdate(tea.Msg) tea.Cmd
	GetFooter() []string
	GetView() string // "" when it shows nothing
	SetSize(size.Size)
}

// activePlugin returns the plugin showing a view, or nil.
func (m *model) activePlugin() plugin {
	for _, p := range m.plugins {
		if p.GetView() != "" {
			return p
		}
	}
	return nil
}

// CapturesInput reports whether a plugin view takes the input, so the side panel must not get it.
func (m *model) CapturesInput() bool {
	return m.activePlugin() != nil
}

// updatePlugins sends msg to every plugin and batches their commands.
func (m *model) updatePlugins(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(m.plugins))
	for _, p := range m.plugins {
		cmds = append(cmds, p.HandleUpdate(msg))
	}
	return tea.Batch(cmds...)
}
