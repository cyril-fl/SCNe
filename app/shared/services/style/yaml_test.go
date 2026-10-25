package style

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/x/ansi"
	plumtransform "github.com/cyril-fl/plum/utils/transform"
	"gopkg.in/yaml.v3"
)

type yamlSample struct {
	Name    string            `yaml:"name"`
	Color   string            `yaml:"color"`
	Empty   string            `yaml:"empty"`
	Example string            `yaml:"example"`
	Items   []yamlSampleItem  `yaml:"items"`
	Tags    []string          `yaml:"tags"`
	None    []string          `yaml:"none"`
	Nested  map[string]string `yaml:"nested"`
}

type yamlSampleItem struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func TestYamlOutputIsValidYAMLOnceUncolored(t *testing.T) {
	want := yamlSample{
		Name:    "scne",
		Color:   "#C9A26B",
		Example: "scne play\nscne config --pretty\n",
		Items:   []yamlSampleItem{{Name: "yaml", Type: "bool"}, {Name: "json", Type: "bool"}},
		Tags:    []string{"a", "b"},
		None:    []string{},
		Nested:  map[string]string{"key": "value"},
	}

	out, err := Yaml(want)
	if err != nil {
		t.Fatalf("Yaml() error = %v", err)
	}

	var got yamlSample
	if err := yaml.Unmarshal([]byte(ansi.Strip(out)), &got); err != nil {
		t.Fatalf("Unmarshal(Yaml()) error = %v\n%s", err, ansi.Strip(out))
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip = %+v, want %+v\n%s", got, want, ansi.Strip(out))
	}
}

func TestYamlMatchesPlumOnceUncolored(t *testing.T) {
	data := yamlSample{
		Name:    "scne",
		Color:   "#C9A26B",
		Example: "scne play\nscne config --pretty\n",
		Items:   []yamlSampleItem{{Name: "yaml", Type: "bool"}},
		Tags:    []string{"a", "b"},
		None:    []string{},
		Nested:  map[string]string{"key": "value"},
	}

	out, err := Yaml(data)
	if err != nil {
		t.Fatalf("Yaml() error = %v", err)
	}

	if got, want := ansi.Strip(out), plumtransform.ParseYaml(data); got != want {
		t.Fatalf("Yaml() uncolored =\n%s\nwant\n%s", got, want)
	}
}
