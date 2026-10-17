package menucodex

func (m *codexMenuItem) GetName() string { return m.id }

type codexMenuItem struct {
	id string
}

// Getters & Setters

// Methods

// Helpers
func makeCodexList() []*codexMenuItem {
	return []*codexMenuItem{
		{id: "codex.locations"},
		{id: "codex.npcs"},
		{id: "codex.beastiary"},
		{id: "codex.objects"},
		{id: "codex.quest"},
	}
}
