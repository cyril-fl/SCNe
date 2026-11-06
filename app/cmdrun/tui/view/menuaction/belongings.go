package menuaction

import (
	"strconv"

	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/components/form"
	"scne/app/shared/components/form/field"
	"scne/app/shared/components/models"
	"scne/app/shared/messages"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
)

const (
	earn  = "earn"
	spend = "spend"
)

// newPickupPage adds a quantity of a codex object to the player: gears go to the bag.
func newPickupPage(gameSession *session.GameSession) *actionPage {
	return newPage(PageTemplate{
		Label: "pickup",
		NewForm: func() *form.Form {
			return form.New(
				field.SelectField(field.SelectTemplate[id.ID]{Template: field.Template[id.ID]{ID: "object", Label: "field.object", Required: true}, Options: getObjectOptions(gameSession)}),
				field.TextField(field.TextTemplate[string]{Template: field.Template[string]{ID: "quantity", Label: "field.quantity", Defaultvalue: "1", Validator: validatePositive, Required: true}}),
			).DisableAutoSubmit().Focus()
		},
		Submit: func(values *form.Form) tea.Cmd {
			entryID, _ := values.GetFieldValue("object").(id.ID)
			qty, _ := strconv.Atoi(values.GetValuesAsMappedString()["quantity"])
			return report(pickup(gameSession, entryID, qty))
		},
	})
}

// newWalletPage earns or spends coins; spending more than the wallet holds fails.
func newWalletPage(gameSession *session.GameSession) *actionPage {
	return newPage(PageTemplate{
		Label: "wallet:apply",
		NewForm: func() *form.Form {
			return form.New(
				field.SelectField(field.SelectTemplate[string]{Template: field.Template[string]{ID: "operation", Label: "field.operation", Required: true}, Options: []models.Item[string]{
					models.NewItem(i19n.Localize("wallet.earn"), "", earn),
					models.NewItem(i19n.Localize("wallet.spend"), "", spend),
				}}),
				coinField("gold"),
				coinField("silver"),
				coinField("copper"),
			).DisableAutoSubmit().Focus()
		},
		Submit: func(values *form.Form) tea.Cmd {
			return report(applyWallet(gameSession, values.GetValuesAsMappedString()))
		},
	})
}

// pickup gives qty of the codex object entryID to the player and returns the chat line.
func pickup(gameSession *session.GameSession, entryID id.ID, qty int) (string, error) {
	player := gameSession.State.Player
	object, s, found := gameSession.State.Codex.Objects.Find(entryID)
	if player == nil || !found {
		return "", logs.Error("error.not_found", map[string]any{"Subject": "item"})
	}

	player.Pick(object, s, qty)
	return i19n.Localize("inventory.msg:picked", map[string]any{"Item": i19n.Localize(object.GetName()), "Qty": qty}), nil
}

// applyWallet earns or spends the coins of values and returns the chat line.
func applyWallet(gameSession *session.GameSession, values map[string]string) (string, error) {
	player := gameSession.State.Player
	if player == nil {
		return "", logs.Error("error.not_found", map[string]any{"Subject": "player"})
	}

	copper := coins(values["gold"])*100 + coins(values["silver"])*10 + coins(values["copper"])

	key := "wallet.msg:earned"
	if values["operation"] == spend {
		if err := player.Wallet.Spend(copper); err != nil {
			return "", err
		}
		key = "wallet.msg:spent"
	} else {
		player.Wallet.Earn(copper)
	}

	w := player.Wallet
	return i19n.Localize(key, map[string]any{"Copper": copper, "Gold": w.Gold, "Silver": w.Silver, "Rest": w.Copper}), nil
}

func coinField(id string) form.Field {
	return field.TextField(field.TextTemplate[string]{Template: field.Template[string]{ID: id, Label: "field." + id, Validator: validateCoins}})
}

// coins reads an empty field as 0.
func coins(value string) int {
	n, _ := strconv.Atoi(value)
	return n
}

// validateCoins accepts empty or a count of at least 0.
func validateCoins(value string) error {
	if n, err := strconv.Atoi(value); value != "" && (err != nil || n < 0) {
		return logs.Error("error.invalid", map[string]any{"Subject": "coins", "Value": value})
	}
	return nil
}

func validatePositive(value string) error {
	if n, err := strconv.Atoi(value); err != nil || n < 1 {
		return logs.Error("error.invalid", map[string]any{"Subject": "quantity", "Value": value})
	}
	return nil
}

// report sends content to the chat, or the error.
func report(content string, err error) tea.Cmd {
	if err != nil {
		return messages.Fail(err)
	}
	return func() tea.Msg { return messages.Result{Content: content} }
}

func getObjectOptions(gameSession *session.GameSession) []models.Item[id.ID] {
	items := make([]models.Item[id.ID], 0)
	for _, entry := range gameSession.State.Codex.Objects.Entries {
		if entry.Object == nil {
			continue
		}
		items = append(items, models.NewItem(i19n.Localize(entry.Object.GetName()), "", entry.ID))
	}
	return items
}
