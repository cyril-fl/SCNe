package primitive

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDescribeAndMutators(t *testing.T) {
	descriptible := Describe(TemplateDescriptible{
		Name:        "name",
		Description: "description",
	})

	if descriptible.GetName() != "name" || descriptible.GetDescription() != "description" {
		t.Fatal("Describe did not initialize the name and description")
	}

	descriptible.SetName("updated name")
	descriptible.SetDescription("updated description")

	if descriptible.GetName() != "updated name" || descriptible.GetDescription() != "updated description" {
		t.Fatal("mutators did not update the name and description")
	}
}

func TestDescriptibleUnmarshalsYAML(t *testing.T) {
	var descriptible Descriptible

	if err := yaml.Unmarshal([]byte("name: yaml name\ndescription: yaml description\n"), &descriptible); err != nil {
		t.Fatalf("yaml.Unmarshal returned an error: %v", err)
	}

	if descriptible.GetName() != "yaml name" {
		t.Fatalf("name = %q, want %q", descriptible.GetName(), "yaml name")
	}
	if descriptible.GetDescription() != "yaml description" {
		t.Fatalf("description = %q, want %q", descriptible.GetDescription(), "yaml description")
	}
}
