package resolvecampaign

import (
	"testing"

	"solopg/app/cmdrun/domain/campaign"
	cmdruntui "solopg/app/cmdrun/tui"
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
