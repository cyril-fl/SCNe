package gameboard

import (
	"solopg/app/shared/messages"
	"solopg/app/shared/services/logs"
	"time"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

func (m *model) handleSaveInput(msg messages.Saved) (*model, tea.Cmd) {
	m.SetErr(msg.Err)

	if m.HasErr() {
		err := logs.Error("error.unexpected:save", map[string]any{
			"Error": msg.Err,
		})

		m.state.Amend(err.Error())
	} else {
		logs.SilentSuccess("campaign:success", map[string]any{
			"Time": time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return m, nil
}

// - Handlers - //
// Input

func (m *model) handleKeyPress(msg tea.KeyPressMsg) (*model, tea.Cmd) {
	switch plumtypes.Key(msg.String()) {
	// case plumtypes.KEY_ENTER:
	// m.handleEnterInput()
	case plumtypes.KEY_UP, plumtypes.KEY_DOWN:
		// m.handleDirectionInput(msg)
		// default:
		// 	handleDefaultInput(m, msg)
	}
	return m.handleCommand(msg)
	// TODO ajouter les autre tyes mes et sourtout gerer qui gere ca avec un tag comme current page
}

func (m *model) handleCommand(msg tea.KeyPressMsg) (*model, tea.Cmd) {
	switch plumtypes.Key(msg.String()) {
	case plumtypes.CMD_CTRL_S:
		return m, messages.Save(m.save)
	}

	return m, nil
}

func (m *model) handleError(msg messages.Failed) (*model, tea.Cmd) {
	m.SetErr(msg.Err)
	return m, nil
}

// - Side Menu - //
// Codex Action
// func (m *model) handleCodexAction(msg codexmenu.Msg) (*model, tea.Cmd) {
// 	_ = msg

// 	return refreshViewport(m, viewoptions.RefreshOption{Positionreset: false})
// }

// Dice Rolled
// func (m *model) handleDiceRolled(msg dicemenu.Msg) (*model, tea.Cmd) {
// 	message := i19n.Localize("dice.roll:result", map[string]any{
// 		"Dice":  msg.Dice,
// 		"Value": msg.Value,
// 	})

// 	_ = message
// 	// -- Advent journal
// 	// TODO  Set l'aureut dans le messahe
// 	// m.state.Amend("Dice", message))

// 	return m, nil
// }

// Hint Rolled
// func (m *model) handleHintRolled(msg hintmenu.Msg) (*model, tea.Cmd) {
// 	localized := make([]string, 0, len(msg.Result))
// 	for _, hint := range msg.Result {
// 		localized = append(localized, i19n.Localize(hint))
// 	}

// 	message := i19n.Localize("hint.roll:result", map[string]any{
// 		"Value": strings.Join(localized, ", "),
// 	})

// 	_ = message

// 	// m.state.Amend("Hint", message))

// 	return m, nil
// }

// Oracle Rolled
// func (m *model) handleOracleRolled(msg oraclemenu.Msg) (*model, tea.Cmd) {
// 	message := i19n.Localize("oracle.roll:result", map[string]any{
// 		"Roll":     msg.Result.Roll,
// 		"Value":    msg.Result.Result,
// 		"Critical": msg.Result.Critical,
// 	})

// 	_ = message

// 	// m.state.Amend("Oracle", message))

// 	return m, nil
// }

// Side Menu Direction
