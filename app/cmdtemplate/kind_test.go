package cmdtemplate

import (
	"bytes"
	"errors"
	"regexp"
	"testing"

	"scne/app/shared/services/strictyaml"
)

// uncomment turns the "# field" lines of a template into YAML, keeping the "##" help lines.
func uncomment(template []byte) []byte {
	return regexp.MustCompile(`(?m)^# ?`).ReplaceAll(template, nil)
}

func TestExistingEntriesAreValid(t *testing.T) {
	for _, k := range kinds() {
		store := k.store()
		names, err := store.Names()
		if err != nil {
			t.Fatalf("%s: Names() error = %v", k.Name, err)
		}
		if len(names) == 0 {
			t.Fatalf("%s: no entry found in %s", k.Name, k.Path)
		}

		for _, name := range names {
			data, err := store.Get(name)
			if err != nil {
				t.Fatalf("%s %q: Get() error = %v", k.Name, name, err)
			}
			if _, err := k.Validate(data); err != nil {
				t.Errorf("%s %q is invalid: %v", k.Name, name, err)
			}
		}
	}
}

func TestTemplatesAreValidOnceNamed(t *testing.T) {
	for _, k := range kinds() {
		template, err := k.Template()
		if err != nil {
			t.Fatalf("%s: Template() error = %v", k.Name, err)
		}

		if _, err := k.Validate(uncomment(template)); err == nil {
			t.Errorf("%s: the unnamed template should be invalid", k.Name)
		}

		if _, err := k.Validate(template); !errors.Is(err, strictyaml.ErrEmpty) {
			t.Errorf("%s: the template should hold only comments, got %v", k.Name, err)
		}

		named := bytes.Replace(uncomment(template), []byte(`name: ""`), []byte(`name: "test"`), 1)
		if name, err := k.Validate(named); err != nil || name != "test" {
			t.Errorf("%s: named template Validate() = %q, %v", k.Name, name, err)
		}
	}
}

func TestValidateRejectsUnknownKeys(t *testing.T) {
	for _, k := range kinds() {
		if _, err := k.Validate([]byte("name: test\nunknown: 1\n")); err == nil {
			t.Errorf("%s: expected an unknown key to be rejected", k.Name)
		}
	}
}
