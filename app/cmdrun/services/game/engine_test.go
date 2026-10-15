package game

import (
	"testing"

	"solopg/app/cmdrun/types/id"
)

func TestNewEngineStoresCampaignAndState(t *testing.T) {
	state := &State{}
	campaignID := id.New()
	engine := NewEngine(campaignID, state)
	if engine.CampaignID != campaignID || engine.State != state {
		t.Fatal("engine did not retain its campaign ID and state")
	}
}
