package hint

import (
	"scne/app/shared/services/logs"
	"scne/app/shared/services/strictyaml"
)

// Validate checks data as one hint file and returns its name.
func Validate(data []byte) (string, error) {
	value, err := strictyaml.Decode[yamlConfig](data)
	if err != nil {
		return "", err
	}

	if value.Name == "" {
		return "", logs.Error("error.required", map[string]any{"Subject": "hint", "Property": "name"})
	}
	if len(value.Values) == 0 {
		return "", logs.Error("error.required", map[string]any{"Subject": value.Name, "Property": "values"})
	}

	return value.Name, nil
}
