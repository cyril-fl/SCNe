package game

import (
	"testing"
	"time"

	"solopg/app/cmdrun/domain/campaign"
	"solopg/app/shared/types/primitive"
)

func TestNewStatePreservesCampaignDates(t *testing.T) {
	campaignCreatedAt := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	campaignUpdatedAt := campaignCreatedAt.Add(time.Hour)
	archiveCreatedAt := campaignCreatedAt.Add(2 * time.Hour)
	archiveUpdatedAt := campaignCreatedAt.Add(3 * time.Hour)

	data := CampaignData{
		Campaign: &campaign.Campaign{
			Datable: primitive.Datable{
				CreatedAt: campaignCreatedAt,
				UpdatedAt: campaignUpdatedAt,
			},
		},
		Archives: &campaign.Archives{
			Datable: primitive.Datable{
				CreatedAt: archiveCreatedAt,
				UpdatedAt: archiveUpdatedAt,
			},
		},
	}

	state := NewState(data)

	if !state.Metadata.Campaign.GetcreatedAt().Equal(campaignCreatedAt) {
		t.Fatal("NewState did not preserve campaign CreatedAt")
	}
	if !state.Metadata.Campaign.GetupdatedAt().Equal(campaignUpdatedAt) {
		t.Fatal("NewState did not preserve campaign UpdatedAt")
	}
	if !state.Metadata.Archive.GetcreatedAt().Equal(archiveCreatedAt) {
		t.Fatal("NewState did not preserve archive CreatedAt")
	}
	if !state.Metadata.Archive.GetupdatedAt().Equal(archiveUpdatedAt) {
		t.Fatal("NewState did not preserve archive UpdatedAt")
	}
}
