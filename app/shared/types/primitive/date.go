package primitive

import "time"

// LOW Je ne suis pas certain quel eBSONC soit necessaire
type Datable struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	immutable bool
}

type TemplateDatable struct {
	Immutable bool
}

func Dated() Datable {
	return newdate(TemplateDatable{
		Immutable: false,
	})
}

func Steady() Datable {
	return newdate(TemplateDatable{
		Immutable: true,
	})
}

func newdate(params TemplateDatable) Datable {
	now := time.Now().UTC()

	return Datable{
		CreatedAt: now,
		UpdatedAt: now,
		immutable: params.Immutable,
	}
}

// Getters & Setters
func (d *Datable) GetcreatedAt() time.Time {
	return d.CreatedAt
}

func (d *Datable) GetupdatedAt() time.Time {
	return d.UpdatedAt
}

func (d *Datable) GetFormattedUpdated(format string) string {
	return d.UpdatedAt.Format(format)
}

func (d *Datable) GetAge() Datable {
	return Datable{
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		immutable: d.immutable,
	}
}

// Methods
func (d *Datable) Touch() {
	if d.immutable {
		return
	}
	d.UpdatedAt = time.Now().UTC()
}
