package quest

import (
	"testing"

	"scne/app/cmdrun/domain/card/attributes/priority"
	"scne/app/cmdrun/domain/card/attributes/variety"
)

func TestNewFallsBackToQuestVarietyWhenInvalid(t *testing.T) {
	for _, provided := range []variety.Variety{"", "mock.quest.unknown"} {
		quest, err := New(Template{Name: "mock.quest", Variety: provided})
		if err != nil {
			t.Fatalf("New(%q) error = %v", provided, err)
		}
		if quest.Variety != "quest_card" {
			t.Errorf("New(%q).Variety = %q, want %q", provided, quest.Variety, "quest_card")
		}
	}
}

func TestNewKeepsValidVarietyAndCopiesFields(t *testing.T) {
	params := Template{
		Name:        "mock.quest",
		Description: "mock.quest.description",
		Variety:     variety.MakeDefault("mock.quest.variety"),
		Done:        true,
		Priority:    *priority.New(priority.Template{Name: "mock.priority", Weight: 2}),
		Rank:        "mock.rank",
		Reward:      "mock.reward",
	}

	quest, err := New(params)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if quest.Variety != params.Variety {
		t.Errorf("Variety = %q, want %q", quest.Variety, params.Variety)
	}
	if quest.Done != params.Done || quest.Priority != params.Priority || quest.Rank != params.Rank || quest.Reward != params.Reward {
		t.Errorf("quest = %#v, want fields from %#v", quest, params)
	}
}
