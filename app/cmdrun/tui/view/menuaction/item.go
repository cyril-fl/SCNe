package menuaction

func (m *actionMenuItem) GetName() string { return m.id }

type actionMenuItem struct {
	id string
}

// Getters & Setters

// Methods

// Helpers
func makeActionList() []*actionMenuItem {
	return []*actionMenuItem{
		{id: FIGHT_PAGE},
		{id: TRAVEL_PAGE},
	}
}
