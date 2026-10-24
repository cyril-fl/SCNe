package codexform

import (
	"fmt"
	"solopg/app/cmdrun/domain/card/attributes/objectcategory"
	"solopg/app/cmdrun/domain/card/attributes/priority"
	"solopg/app/cmdrun/domain/card/attributes/rank"
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
	"strconv"
)

// - NPC Form - //
func GetNpcForm(values map[string]string) *form.Form {
	statsFields := getFormStatsFields(values, "")
	encounter := values["encounter"]
	if encounter == "" {
		encounter = "encounter_normal"
	}
	return form.
		New(
			field.TextField(field.TextTemplate[string]{ID: "name", Label: "field.name", Defaultvalue: values["name"], Validator: nil, Required: true}),
			field.TextField(field.TextTemplate[string]{ID: "description", Label: "field.description", Defaultvalue: values["description"], Validator: nil, Required: true}),
			field.SelectField(field.SelectTemplate[string]{ID: "class", Label: "field.class", Defaultvalue: values["class"], Options: getFormClassOptions(), Validator: nil, Required: true}),
			field.SelectField(field.SelectTemplate[string]{ID: "race", Label: "field.race", Defaultvalue: values["race"], Options: getFormRaceOptions(), Validator: nil, Required: true}),
			field.SelectField(field.SelectTemplate[string]{ID: "encounter", Label: "field.encounter", Defaultvalue: encounter, Options: getFormOracle(), Validator: nil, Required: true}),
		).
		Add(statsFields...).
		Focus()
}

// - Beast Form - //
func GetBeastiaryForm(values map[string]string) *form.Form {
	statsFields := getFormStatsFields(values, "")
	encounter := values["encounter"]
	if encounter == "" {
		encounter = "encounter_normal"
	}
	return form.
		New(
			field.TextField(field.TextTemplate[string]{ID: "name", Label: "field.name", Defaultvalue: values["name"], Validator: nil, Required: true}),
			field.TextField(field.TextTemplate[string]{ID: "description", Label: "field.description", Defaultvalue: values["description"], Validator: nil, Required: true}),
			field.SelectField(field.SelectTemplate[string]{ID: "encounter", Label: "field.encounter", Defaultvalue: encounter, Options: getFormOracle(), Validator: nil, Required: true}),
		).
		Add(statsFields...).
		Focus()
}

// - Location Form - //
func GetLocationForm(values map[string]string) *form.Form {
	return form.New(
		field.TextField(field.TextTemplate[string]{ID: "name", Label: "field.name", Defaultvalue: values["name"], Validator: nil, Required: true}),
		field.TextField(field.TextTemplate[string]{ID: "description", Label: "field.description", Defaultvalue: values["description"], Validator: nil, Required: true}),
		form.Repeat(form.RepeatTemplate{
			ID:    "effect",
			Count: countEffects(values),
			Make: func(index int) []form.Field {
				return getEffectFields(values, index)
			},
		}),
	).Focus()
}

func countEffects(values map[string]string) int {
	count := 0
	for {
		if _, ok := values[fmt.Sprintf("effect.%d.name", count)]; !ok {
			return count
		}
		count++
	}
}

func getEffectFields(values map[string]string, index int) []form.Field {
	prefix := fmt.Sprintf("effect.%d.", index)
	fields := []form.Field{
		field.TextField(field.TextTemplate[string]{ID: prefix + "name", Label: "field.effect.name", Defaultvalue: values[prefix+"name"], Validator: nil, Required: false}),
		field.TextField(field.TextTemplate[string]{ID: prefix + "description", Label: "field.effect.description", Defaultvalue: values[prefix+"description"], Validator: nil, Required: false}),
	}
	return append(fields, getFormStatsFields(values, prefix)...)
}

// - Object Form - //
func GetObjectForm(values map[string]string) *form.Form {
	var categoryItems []models.Item[string]

	if len(objectcategory.List()) == 0 {
		objectcategory.MakeDefault("object")
	}

	for _, i := range objectcategory.List() {
		categoryItems = append(categoryItems, models.NewItem(i19n.Localize(i.String()), "", i.String()))
	}

	return form.New(
		field.TextField(field.TextTemplate[string]{ID: "name", Label: "field.name", Defaultvalue: values["name"], Validator: nil, Required: true}),
		field.TextField(field.TextTemplate[string]{ID: "description", Label: "field.description", Defaultvalue: values["description"], Validator: nil, Required: true}),
		field.SelectField(field.SelectTemplate[string]{ID: "category", Label: "field.category", Defaultvalue: values["category"], Options: categoryItems, Validator: nil, Required: true}),
	).Focus()
}

// - Objectif Form - //
func GetQuestForm(values map[string]string) *form.Form {
	if len(priority.List()) == 0 {
		priority.MakeDefault("priority.medium")
	}

	if len(rank.List()) == 0 {
		rank.MakeDefault("side")
	}

	return form.New(
		field.TextField(field.TextTemplate[string]{ID: "name", Label: "field.name", Defaultvalue: values["name"], Validator: nil, Required: true}),
		field.TextField(field.TextTemplate[string]{ID: "description", Label: "field.description", Defaultvalue: values["description"], Validator: nil, Required: true}),
		field.SelectField(field.SelectTemplate[string]{ID: "done", Label: "field.done", Defaultvalue: values["done"], Options: getDoneOptions(), Validator: nil, Required: false}),
		field.TextField(field.TextTemplate[string]{ID: "reward", Label: "field.reward", Defaultvalue: values["reward"], Validator: nil, Required: false}),
		field.SelectField(field.SelectTemplate[string]{ID: "priority", Label: "field.priority", Defaultvalue: values["priority"], Options: getFormPriorityOptions(), Validator: nil, Required: false}),
		field.SelectField(field.SelectTemplate[string]{ID: "rank", Label: "field.rank", Defaultvalue: values["rank"], Options: getFormRankOptions(), Validator: nil, Required: false}),
	).Focus()
}

func getDoneOptions() []models.Item[string] {
	return []models.Item[string]{
		models.NewItem(i19n.Localize("field.done:true"), "", "true"),
		models.NewItem(i19n.Localize("field.done:false"), "", "false"),
	}
}

func getFormPriorityOptions() []models.Item[string] {
	values, err := priority.ListWithError()
	if err != nil {
		return nil
	}
	items := make([]models.Item[string], 0, len(values))
	for _, value := range values {
		items = append(items, models.NewItem(i19n.Localize(value.GetName()), "", value.GetName()))
	}
	return items
}

func getFormRankOptions() []models.Item[string] {
	values, err := rank.ListWithError()
	if err != nil {
		return nil
	}
	items := make([]models.Item[string], 0, len(values))
	for _, value := range values {
		items = append(items, models.NewItem(i19n.Localize(value.String()), "", value.String()))
	}
	return items
}

// - Helpers - //

// Helpers
func getFormOracle() []models.Item[string] {
	values, err := oracle.ListFromFolder(config.GetConfig().Documents.Folders.Encounters)
	if err != nil {
		return nil
	}
	items := make([]models.Item[string], 0, len(values))
	for _, value := range values {
		items = append(items, models.NewItem(i19n.Localize("oracles."+value), "", value))
	}
	return items
}

func getFormStatsFields(values map[string]string, prefix string) []form.Field {
	fields := make([]form.Field, 0)
	for _, value := range stats.List() {
		fields = append(fields, field.TextField(field.TextTemplate[string]{
			ID:           prefix + value.String(),
			Label:        i19n.Localize(value.String()),
			Defaultvalue: values[prefix+value.String()],
			Validator:    validateStat,
		}))
	}
	return fields
}

// Empty is allowed: the stat is then rolled from the encounter oracle.
func validateStat(value string) error {
	if value == "" {
		return nil
	}
	if _, err := strconv.Atoi(value); err != nil {
		return logs.Error("error.invalid", map[string]any{
			"Subject": "stat",
			"Value":   value,
		})
	}
	return nil
}

func getFormRaceOptions() []models.Item[string] {
	items := make([]models.Item[string], 0)
	for _, value := range races.List() {
		items = append(items, models.NewItem(i19n.Localize(value.GetName()), "", value.GetID()))
	}
	return items
}

func getFormClassOptions() []models.Item[string] {
	items := make([]models.Item[string], 0)
	for _, value := range classes.List() {
		items = append(items, models.NewItem(i19n.Localize(value.GetName()), "", value.GetName()))
	}
	return append(items, models.NewItem("None", "", ""))
}
