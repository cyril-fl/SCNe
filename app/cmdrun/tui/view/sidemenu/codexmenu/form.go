package codexmenu

import (
	"solopg/app/cmdrun/domain/card/attributes/objectcategory"
	"solopg/app/shared/components/form"
	"solopg/app/shared/components/form/field"
	"solopg/app/shared/components/models"
	"solopg/app/shared/services/i19n"
)

// - NPC Form - //
func getNpcForm(values map[string]string) *form.Form {
	statsFields := getFormStatsFields(values)
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
func getBeastiaryForm(values map[string]string) *form.Form {
	statsFields := getFormStatsFields(values)
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
func getLocationForm(values map[string]string) *form.Form {
	return form.New(
		field.TextField(field.TextTemplate[string]{ID: "name", Label: "field.name", Defaultvalue: values["name"], Validator: nil, Required: true}),
		field.TextField(field.TextTemplate[string]{ID: "description", Label: "field.description", Defaultvalue: values["description"], Validator: nil, Required: true}),
	).Focus()
}

// - Object Form - //
func getObjectForm(values map[string]string) *form.Form {
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
func getObjectifForm(values map[string]string) *form.Form {
	return form.New(
		field.TextField(field.TextTemplate[string]{ID: "name", Label: "field.name", Defaultvalue: values["name"], Validator: nil, Required: true}),
		field.TextField(field.TextTemplate[string]{ID: "description", Label: "field.description", Defaultvalue: values["description"], Validator: nil, Required: true}),
	).Focus()
}
