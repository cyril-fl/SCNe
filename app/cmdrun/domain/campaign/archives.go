package campaign

import (
	"solopg/app/cmdrun/domain/gameplay/codex"
	"solopg/app/cmdrun/types/id"
	"solopg/app/cmdrun/types/interfaces"
	"solopg/app/shared/types/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

// TEST
type Archives struct {
	primitive.Datable

	CampaignID id.ID        `bson:"campaignId" json:"campaignId"`
	Codex      *codex.Codex `bson:"codex" json:"codex"`
	Journal    *Journal     `bson:"journal" json:"journal"`
	Log        *Journal     `bson:"log" json:"log"`
}

type ArchivesTemplate struct {
	CampaignID id.ID
	Codex      *codex.Codex
	Journal    *Journal
	Log        *Journal
}

func NewArchives(params ArchivesTemplate) *Archives {
	return &Archives{
		Datable: primitive.Dated(),

		CampaignID: params.CampaignID,
		Codex:      params.Codex,
		Journal:    params.Journal,
		Log:        params.Log,
	}
}

func (a *Archives) Filter() bson.M {
	return bson.M{
		"campaignId": a.CampaignID,
	}
}

func (a *Archives) EnsureInitialized() {
	if a.Codex == nil {
		a.Codex = codex.New()
	}

	if a.Journal == nil {
		a.Journal = NewJournal([]Entry{})
	}

	if a.Log == nil {
		a.Log = NewJournal([]Entry{})
	}

	list := []interfaces.Initializable{
		a.Codex,
		a.Journal,
		a.Log,
	}

	for _, archive := range list {
		archive.EnsureInitialized()
	}
}
