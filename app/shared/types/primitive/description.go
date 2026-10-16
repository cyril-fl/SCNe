package primitive

// TEST

type Descriptible struct {
	Name        string
	Description string
}

type TemplateDescriptible struct {
	Name        string //`bson:"name" json:"name" yaml:"name"`
	Description string //`bson:"description" json:"description" yaml:"description"`
}

func Describe(params TemplateDescriptible) Descriptible {
	return Descriptible{
		Name:        params.Name,
		Description: params.Description,
	}
}

// Getters & Setters
func (d Descriptible) GetName() string {
	return d.Name
}

func (d Descriptible) GetDescription() string {
	return d.Description
}

func (d *Descriptible) SetName(name string) {
	d.Name = name
}

func (d *Descriptible) SetDescription(description string) {
	d.Description = description
}
