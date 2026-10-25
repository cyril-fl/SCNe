package campaign

import (
	"solopg/app/cmdrun/domain/gameplay/codex"
	"solopg/app/cmdrun/types/id"
	"solopg/app/cmdrun/types/interfaces"

	plumprimitive "github.com/cyril-fl/plum/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

// Archives holds the codex, journal and log of a campaign.
type Archives struct {
	plumprimitive.Datable

	CampaignID id.ID        `bson:"campaignId" json:"campaignId"`
	Codex      *codex.Codex `bson:"codex" json:"codex"`
	Journal    *Journal     `bson:"journal" json:"journal"`
	Log        *Journal     `bson:"log" json:"log"`
}

// ArchivesTemplate configures Archives.
type ArchivesTemplate struct {
	CampaignID id.ID
	Codex      *codex.Codex
	Journal    *Journal
	Log        *Journal
}

// NewArchives returns dated Archives.
func NewArchives(params ArchivesTemplate) *Archives {
	return &Archives{
		Datable: plumprimitive.Dated(),

		CampaignID: params.CampaignID,
		Codex:      params.Codex,
		Journal:    params.Journal,
		Log:        params.Log,
	}
}

// Filter returns the MongoDB filter matching the campaign archives.
func (a *Archives) Filter() bson.M {
	return bson.M{
		"campaignId": a.CampaignID,
	}
}

// EnsureInitialized creates the missing codex, journal and log.
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
