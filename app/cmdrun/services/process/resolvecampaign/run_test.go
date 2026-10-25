package resolvecampaign

import (
	"testing"

	"scne/app/cmdrun/domain/campaign"
	"scne/app/cmdrun/domain/card/characters/races"
	cmdruntui "scne/app/cmdrun/tui"
)

func TestRunIdentifiesNewCampaign(t *testing.T) {
	if !Process(&cmdruntui.Context{}).isNewCampaign() {
		t.Fatal("empty context should identify a new campaign")
	}
}

func TestRunKeepsSelectedCampaign(t *testing.T) {
	selected := &campaign.Campaign{}
	p := Process(&cmdruntui.Context{SelectedSave: selected})
	p.Run()

	if p.GetResult() != selected {
		t.Fatal("selected campaign was not returned")
	}
	if p.HasErr() {
		t.Fatalf("unexpected error: %v", p.GetErr())
	}
}

func TestRunRejectsIncompleteOnboarding(t *testing.T) {
	p := Process(&cmdruntui.Context{SelectedRace: &races.Race{}})
	p.Run()

	if !p.HasErr() {
		t.Fatal("missing class and location should fail")
	}
	if p.GetResult() != nil {
		t.Fatal("no campaign should be built from an incomplete onboarding")
	}
}
