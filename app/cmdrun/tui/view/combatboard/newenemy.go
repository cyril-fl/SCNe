package combatboard

import (
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/cmdrun/tui/shortcut"
	"scne/app/cmdrun/tui/view/menuinventory"
	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/form"
	"scne/app/shared/messages"

	codexform "scne/app/cmdrun/domain/gameplay/codex/form"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

// formPlugin opens a form on key and runs submit once the form is valid.
type formPlugin struct {
	key      plumtypes.Key
	shortcut shortcut.Shortcut
	newForm  func() *form.Form
	submit   func(values *form.Form) tea.Cmd

	form *form.Form
	size size.Size
}

// newEnemyPlugin opens the codex NPC form on ctrl+n; the saved NPC joins the codex and the combat.
func newEnemyPlugin(gameSession *session.GameSession, onAdd func()) *formPlugin {
	return &formPlugin{
		key:      plumtypes.SHORT_CTRL_N,
		shortcut: shortcut.NEW,
		newForm:  func() *form.Form { return codexform.GetNpcForm(map[string]string{}) },
		submit: func(values *form.Form) tea.Cmd {
			// The codex rolls the stats.
			npcs := gameSession.State.Codex.Npcs
			if err := npcs.AddFromMap(values.GetValuesAsMappedString()); err != nil {
				return messages.Fail(err)
			}
			gameSession.AddOpponent(npcs.Entries[len(npcs.Entries)-1].Character)
			onAdd()
			return nil
		},
	}
}

// newInventoryPlugin opens the inventory form on ctrl+o.
func newInventoryPlugin(gameSession *session.GameSession) *formPlugin {
	return &formPlugin{
		key:      plumtypes.SHORT_CTRL_O,
		shortcut: shortcut.INVENTORY,
		newForm:  func() *form.Form { return menuinventory.NewForm(gameSession) },
		submit: func(values *form.Form) tea.Cmd {
			return menuinventory.Submit(gameSession, values)
		},
	}
}

// newVitalsPlugin opens the health form on ctrl+d, to note damage taken or healing.
func newVitalsPlugin(gameSession *session.GameSession) *formPlugin {
	return &formPlugin{
		key:      plumtypes.SHORT_CTRL_D,
		shortcut: shortcut.VITALS,
		newForm:  menuinventory.NewVitalsForm,
		submit: func(values *form.Form) tea.Cmd {
			return menuinventory.SubmitVitals(gameSession, values)
		},
	}
}

// Getters & Setters

func (p *formPlugin) GetView() string {
	if p.form == nil {
		return ""
	}
	return p.form.View().Content
}

func (p *formPlugin) GetFooter() []string {
	if p.form == nil {
		return []string{p.shortcut.String()}
	}
	return []string{shortcut.SAVEFORM.String(), shortcut.LEAVE.String()}
}

func (p *formPlugin) SetSize(s size.Size) {
	p.size = s
	if p.form != nil {
		p.form.SetWidth(s.GetWidth()).SetHeight(s.GetHeight())
	}
}

// Methods

func (p *formPlugin) HandleUpdate(msg tea.Msg) tea.Cmd {
	if p.form == nil {
		if key, ok := msg.(tea.KeyPressMsg); ok && plumtypes.Key(key.String()) == p.key {
			p.form = p.newForm()
			p.SetSize(p.size)
		}
		return nil
	}

	switch msg := msg.(type) {
	case form.Validate:
		return p.handleValidate()
	case tea.KeyPressMsg:
		switch plumtypes.Key(msg.String()) {
		case plumtypes.KEY_ESC:
			p.form = nil
			return nil
		case plumtypes.CMD_CTRL_ENTER:
			return form.Send[form.Validate]()
		}
	}

	_, cmd := p.form.Update(msg)
	return cmd
}

// handleValidate keeps the form open on errors and closes it after submit.
func (p *formPlugin) handleValidate() tea.Cmd {
	p.form.Submit()
	if p.form.HasErr() {
		return messages.Fail(p.form.GetErr())
	}

	cmd := p.submit(p.form)
	p.form = nil
	return cmd
}
