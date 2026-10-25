package campaign

import (
	"scne/app/cmdrun/domain/gameplay/codex"
	"scne/app/cmdrun/types/id"
	"scne/app/cmdrun/types/interfaces"

	plumprimitive "github.com/cyril-fl/plum/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

// Archives holds the codex, journal and log of a campaign.
type Archives struct {
	plumprimitive.Datable

	CampaignID id.ID        `bson:"campaignId" json:"campaignId"`
	Codex      *codex.Codex `bson:"codex" json:"codex"`
	Journal    *Journal     `bson:"journal" json:"journal"`
}

// ArchivesTemplate configures Archives.
type ArchivesTemplate struct {
	CampaignID id.ID
	Codex      *codex.Codex
	Journal    *Journal
}

// NewArchives returns dated Archives.
func NewArchives(params ArchivesTemplate) *Archives {
	return &Archives{
		Datable: plumprimitive.Dated(),

		CampaignID: params.CampaignID,
		Codex:      params.Codex,
		Journal:    params.Journal,
	}
}

// Filter returns the MongoDB filter matching the campaign archives.
func (a *Archives) Filter() bson.M {
	return bson.M{
		"campaignId": a.CampaignID,
	}
}

// EnsureInitialized creates the missing codex and journal.
func (a *Archives) EnsureInitialized() {
	if a.Codex == nil {
		a.Codex = codex.New()
	}

	if a.Journal == nil {
		a.Journal = NewJournal([]Entry{})
	}

	list := []interfaces.Initializable{
		a.Codex,
		a.Journal,
	}

	for _, archive := range list {
		archive.EnsureInitialized()
	}
}
