package strictyaml

import (
	"bytes"
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

// ErrEmpty is returned for a document holding nothing but comments.
var ErrEmpty = errors.New("empty document")

// Decode decodes data into T, rejecting the keys T does not declare.
func Decode[T any](data []byte) (T, error) {
	var value T

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&value); err != nil {
		if errors.Is(err, io.EOF) {
			return value, ErrEmpty
		}
		return value, err
	}

	return value, nil
}
