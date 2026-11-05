package characters

import (
	"scne/app/cmdrun/domain/card/attributes/rarity"
	"scne/app/cmdrun/domain/card/attributes/variety"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/strictyaml"
)

// Validate checks data as one character card and returns its name.
// The class is not checked: loading adds a missing class, see newFromYaml.
func Validate(data []byte) (string, error) {
	value, err := strictyaml.Decode[yamlConfig](data)
	if err != nil {
		return "", err
	}

	if value.Name == "" {
		return "", logs.Error("error.required", map[string]any{"Subject": "character", "Property": "name"})
	}

	_, err = New(Template{
		Name:        value.Name,
		Description: value.Description,
		Race:        value.Race,
		Rarity:      rarity.Rarity(value.Rarity),
		Variety:     variety.Variety(value.Variety),
		Stats:       value.Stats,
		Wallet:      value.Wallet,
		Equipment:   value.Equipment,
		Inventory:   value.Inventory,
	})
	return value.Name, err
}
