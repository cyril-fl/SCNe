package races

import (
	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/strictyaml"
)

// Validate checks data as one entry of the races file and returns its name.
func Validate(data []byte) (string, error) {
	value, err := strictyaml.Decode[yamlConfig](data)
	if err != nil {
		return "", err
	}

	if value.Name == "" {
		return "", logs.Error("error.required", map[string]any{"Subject": "race", "Property": "name"})
	}

	return value.Name, stats.ValidateModifiers(value.Bonus)
}
