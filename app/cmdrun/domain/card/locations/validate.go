package locations

import (
	"scne/app/cmdrun/domain/card/attributes/rarity"
	"scne/app/cmdrun/domain/card/attributes/variety"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/strictyaml"
)

// Validate checks data as one location card and returns its name.
func Validate(data []byte) (string, error) {
	value, err := strictyaml.Decode[yamlConfig](data)
	if err != nil {
		return "", err
	}

	if value.Name == "" {
		return "", logs.Error("error.required", map[string]any{"Subject": "location", "Property": "name"})
	}

	_, err = New(Template{
		Name:        value.Name,
		Description: value.Description,
		Rarity:      rarity.Rarity(value.Rarity),
		Variety:     variety.Variety(value.Variety),
		Effects:     value.Effects,
	})
	return value.Name, err
}
