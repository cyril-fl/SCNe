package style

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Json renders data as indented JSON colored like Yaml; uncolored, it matches json.MarshalIndent with two spaces.
func Json(data any) (string, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Tokens keep the struct field order, which decoding into a map would lose.
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()

	var b strings.Builder
	if err := writeJsonValue(&b, decoder, 0); err != nil {
		return "", err
	}

	return b.String(), nil
}

// Helpers

func writeJsonValue(b *strings.Builder, decoder *json.Decoder, indent int) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	switch value := token.(type) {
	case json.Delim:
		return writeJsonContainer(b, decoder, indent, value)
	case string:
		if value == "" {
			b.WriteString(Muted().Render(`""`))
			return nil
		}
		quoted, _ := json.Marshal(value)
		b.WriteString(text().Render(string(quoted)))
	case nil:
		b.WriteString(Muted().Render("null"))
	default:
		b.WriteString(text().Render(fmt.Sprint(value)))
	}

	return nil
}

// writeJsonContainer writes an object or an array whose opening delimiter was already read.
func writeJsonContainer(b *strings.Builder, decoder *json.Decoder, indent int, open json.Delim) error {
	isObject := open == '{'
	closing := "]"
	if isObject {
		closing = "}"
	}

	if !decoder.More() {
		_, err := decoder.Token()
		b.WriteString(Muted().Render(open.String() + closing))
		return err
	}

	pad := strings.Repeat(" ", indent+2)
	b.WriteString(Muted().Render(open.String()) + "\n")

	for decoder.More() {
		b.WriteString(pad)
		if isObject {
			key, err := decoder.Token()
			if err != nil {
				return err
			}
			quoted, _ := json.Marshal(key)
			b.WriteString(Accent().Render(string(quoted)) + Muted().Render(":") + " ")
		}
		if err := writeJsonValue(b, decoder, indent+2); err != nil {
			return err
		}
		if decoder.More() {
			b.WriteString(Muted().Render(","))
		}
		b.WriteString("\n")
	}

	if _, err := decoder.Token(); err != nil {
		return err
	}
	b.WriteString(strings.Repeat(" ", indent) + Muted().Render(closing))

	return nil
}
