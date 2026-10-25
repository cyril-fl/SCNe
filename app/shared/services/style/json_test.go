package style

import (
	"encoding/json"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestJsonMatchesMarshalIndentOnceUncolored(t *testing.T) {
	data := yamlSample{
		Name:    "scne",
		Color:   "#C9A26B",
		Example: "scne play\n<tag> & \"quotes\"\n",
		Items:   []yamlSampleItem{{Name: "yaml", Type: "bool"}, {Name: "json", Type: "bool"}},
		Tags:    []string{"a", "b"},
		None:    []string{},
		Nested:  map[string]string{"b": "2", "a": "1"},
	}

	out, err := Json(data)
	if err != nil {
		t.Fatalf("Json() error = %v", err)
	}

	want, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if got := ansi.Strip(out); got != string(want) {
		t.Fatalf("Json() uncolored =\n%s\nwant\n%s", got, want)
	}
}
