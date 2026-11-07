package menuinventory

import (
	"strconv"

	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/shared/components/form"
	"scne/app/shared/components/form/field"
	"scne/app/shared/components/models"
	"scne/app/shared/messages"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/logs"

	tea "charm.land/bubbletea/v2"
	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

const (
	lose    = "lose"
	recover = "recover"
)

// NewVitalsForm lets the player lose or recover health or stamina: the amount is theirs to decide.
func NewVitalsForm() *form.Form {
	resources := []models.Item[string]{}
	for _, stat := range stats.List() {
		if characters.IsResource(stat) {
			resources = append(resources, models.NewItem(i19n.Localize("stat."+stat.String()), "", stat.String()))
		}
	}

	return form.New(
		field.SelectField(field.SelectTemplate[string]{Template: field.Template[string]{ID: "stat", Label: "field.stat", Required: true}, Options: resources}),
		field.SelectField(field.SelectTemplate[string]{Template: field.Template[string]{ID: "operation", Label: "field.operation", Required: true}, Options: []models.Item[string]{
			models.NewItem(i19n.Localize("vitals.lose"), "", lose),
			models.NewItem(i19n.Localize("vitals.recover"), "", recover),
		}}),
		field.TextField(field.TextTemplate[string]{Template: field.Template[string]{ID: "amount", Label: "field.amount", Validator: validateAmount, Required: true}}),
	).DisableAutoSubmit().Focus()
}

// SubmitVitals applies the form to the player and reports the new value in the chat.
func SubmitVitals(gameSession *session.GameSession, values *form.Form) tea.Cmd {
	result, err := gameSession.Track(func() (string, error) {
		return applyVitals(gameSession.State.Player, values.GetValuesAsMappedString())
	})
	if err != nil {
		return messages.Fail(err)
	}
	return func() tea.Msg { return messages.Result{Content: result} }
}

func applyVitals(player *characters.Character, values map[string]string) (string, error) {
	if player == nil {
		return "", logs.Error("error.not_found", map[string]any{"Subject": "player"})
	}

	stat := stats.Stat(values["stat"])
	amount, _ := strconv.Atoi(values["amount"])

	key := "vitals.msg:lost"
	if values["operation"] == recover {
		player.Heal(stat, amount)
		key = "vitals.msg:recovered"
	} else {
		player.Damage(stat, amount)
	}

	return i19n.Localize(key, map[string]any{
		"Stat":   plumtransform.Capitalize(i19n.Localize("stat." + stat.String())),
		"Amount": amount,
		"Value":  player.FormatStat(stat),
	}), nil
}

func validateAmount(value string) error {
	if n, err := strconv.Atoi(value); err != nil || n < 1 {
		return logs.Error("error.invalid", map[string]any{"Subject": "amount", "Value": value})
	}
	return nil
}
