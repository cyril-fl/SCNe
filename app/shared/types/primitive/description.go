package primitive

import "gopkg.in/yaml.v3"

// TEST

type Descriptible struct {
	name        string
	description string
}

type TemplateDescriptible struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func Describe(params TemplateDescriptible) Descriptible {
	return Descriptible{
		name:        params.Name,
		description: params.Description,
	}
}

// Getters & Setters
func (d Descriptible) GetName() string {
	return d.name
}

func (d Descriptible) GetDescription() string {
	return d.description
}

func (d *Descriptible) SetName(name string) {
	d.name = name
}

func (d *Descriptible) SetDescription(description string) {
	d.description = description
}

func (d *Descriptible) UnmarshalYAML(value *yaml.Node) error {
	var template TemplateDescriptible
	if err := value.Decode(&template); err != nil {
		return err
	}

	d.name = template.Name
	d.description = template.Description
	return nil
}
