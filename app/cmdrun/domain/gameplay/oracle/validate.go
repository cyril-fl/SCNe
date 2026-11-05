package oracle

import (
	"scne/app/shared/services/logs"
	"scne/app/shared/services/strictyaml"
)

// Validate checks data as one oracle file and returns its name.
func Validate(data []byte) (string, error) {
	value, err := strictyaml.Decode[yamlConfig](data)
	if err != nil {
		return "", err
	}

	if value.Name == "" {
		return "", logs.Error("error.required", map[string]any{"Subject": "oracle", "Property": "name"})
	}

	return value.Name, check(value)
}
