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
		{id: INVENTORY_PAGE},
		{id: PICKUP_PAGE},
		{id: WALLET_PAGE},
		{id: VITALS_PAGE},
		{id: CONDITIONS_PAGE},
	}
}
