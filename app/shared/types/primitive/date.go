package primitive

import "time"

// TEST

type Datable struct {
	createdAt time.Time `bson:"createdAt" json:"createdAt"`
	updatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
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
		createdAt: now,
		updatedAt: now,
		immutable: params.Immutable,
	}
}

// Getters & Setters
func (d *Datable) GetcreatedAt() time.Time {
	return d.createdAt
}

func (d *Datable) GetupdatedAt() time.Time {
	return d.updatedAt
}


func (d *Datable) GetFormattedUpdated(format string) string {
	return d.updatedAt.Format(format)
}


func (d *Datable) GetAge() Datable {
	return Datable{
		createdAt: d.createdAt,
		updatedAt: d.updatedAt,
	}
}

// Methods
func (d *Datable) Touch() {
	if d.immutable {
		return
	}
	d.updatedAt = time.Now().UTC()
}


