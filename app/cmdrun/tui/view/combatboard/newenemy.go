package combatboard

import (
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/cmdrun/tui/shortcut"
	"scne/app/cmdrun/types/size"
	"scne/app/shared/components/form"
	"scne/app/shared/messages"

	codexform "scne/app/cmdrun/domain/gameplay/codex/form"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

// newEnemy opens the codex NPC form on ctrl+n; the saved NPC joins the codex and the combat.
type newEnemy struct {
	session *session.GameSession
	form    *form.Form
	size    size.Size
	onAdd   func()
}

func newEnemyPlugin(gameSession *session.GameSession, onAdd func()) *newEnemy {
	return &newEnemy{session: gameSession, onAdd: onAdd}
}

// Getters & Setters

func (p *newEnemy) GetView() string {
	if p.form == nil {
		return ""
	}
	return p.form.View().Content
}

func (p *newEnemy) GetFooter() []string {
	if p.form == nil {
		return []string{shortcut.NEW.String()}
	}
	return []string{shortcut.SAVEFORM.String(), shortcut.LEAVE.String()}
}

func (p *newEnemy) SetSize(s size.Size) {
	p.size = s
	if p.form != nil {
		p.form.SetWidth(s.GetWidth()).SetHeight(s.GetHeight())
	}
}

// Methods

func (p *newEnemy) HandleUpdate(msg tea.Msg) tea.Cmd {
	if p.form == nil {
		if key, ok := msg.(tea.KeyPressMsg); ok && plumtypes.Key(key.String()) == plumtypes.SHORT_CTRL_N {
			p.form = codexform.GetNpcForm(map[string]string{})
			p.SetSize(p.size)
		}
		return nil
	}

	switch msg := msg.(type) {
	case form.Validate:
		return p.submit()
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

// submit saves the NPC in the codex, where its stats are rolled, and adds it to the combat.
func (p *newEnemy) submit() tea.Cmd {
	p.form.Submit()
	if p.form.HasErr() {
		return messages.Fail(p.form.GetErr())
	}

	npcs := p.session.State.Codex.Npcs
	if err := npcs.AddFromMap(p.form.GetValuesAsMappedString()); err != nil {
		return messages.Fail(err)
	}

	p.session.AddOpponent(npcs.Entries[len(npcs.Entries)-1].Character)
	p.form = nil
	p.onAdd()
	return nil
}
