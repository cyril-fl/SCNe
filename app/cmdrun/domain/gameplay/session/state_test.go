package session

import (
	"testing"
	"time"

	"solopg/app/cmdrun/domain/campaign"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

func TestNewSessionStatePreservesCampaignDates(t *testing.T) {
	campaignCreatedAt := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	campaignUpdatedAt := campaignCreatedAt.Add(time.Hour)
	archiveCreatedAt := campaignCreatedAt.Add(2 * time.Hour)
	archiveUpdatedAt := campaignCreatedAt.Add(3 * time.Hour)

	data := SessionData{
		Campaign: &campaign.Campaign{
			Datable: plumprimitive.Datable{
				CreatedAt: campaignCreatedAt,
				UpdatedAt: campaignUpdatedAt,
			},
		},
		Archives: &campaign.Archives{
			Datable: plumprimitive.Datable{
				CreatedAt: archiveCreatedAt,
				UpdatedAt: archiveUpdatedAt,
			},
		},
	}

	state := NewSessionState(data)

	if !state.Metadata.Campaign.GetcreatedAt().Equal(campaignCreatedAt) {
		t.Fatal("NewSessionState did not preserve campaign CreatedAt")
	}
	if !state.Metadata.Campaign.GetupdatedAt().Equal(campaignUpdatedAt) {
		t.Fatal("NewSessionState did not preserve campaign UpdatedAt")
	}
	if !state.Metadata.Archive.GetcreatedAt().Equal(archiveCreatedAt) {
		t.Fatal("NewSessionState did not preserve archive CreatedAt")
	}
	if !state.Metadata.Archive.GetupdatedAt().Equal(archiveUpdatedAt) {
		t.Fatal("NewSessionState did not preserve archive UpdatedAt")
	}
}
