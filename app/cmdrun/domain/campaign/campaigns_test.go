package campaign

import (
	"testing"

	"scne/app/cmdrun/domain/card"
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/card/locations"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/i19n"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

func TestNewCampaign(t *testing.T) {
	player := &characters.Character{}
	location := &locations.Location{}
	campaign := New(Template{Player: player, CurrentLocation: location})

	if campaign.ID == (id.ID{}) || campaign.Player != player || campaign.CurrentLocation != location {
		t.Fatal("New did not initialize or copy the campaign")
	}
}

func TestCampaignFilter(t *testing.T) {
	campaign := New(Template{})

	if got := campaign.Filter()["Id"]; got != campaign.ID {
		t.Fatalf("Filter Id = %v, want %v", got, campaign.ID)
	}
}

func TestCampaignTitle(t *testing.T) {
	i19n.SetCache(nil)
	player := &characters.Character{
		Card: card.Card{Descriptible: plumprimitive.Describe(plumprimitive.TemplateDescriptible{Name: "Aria"})},
		Race: "elf", Class: "warrior",
	}

	campaign := New(Template{Player: player})
	if got, want := campaign.Title(), "Aria | {{elf}} | {{warrior}}"; got != want {
		t.Fatalf("Title() = %q, want %q", got, want)
	}
}
