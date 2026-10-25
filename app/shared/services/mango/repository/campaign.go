package repository

import (
	"scne/app/cmdrun/domain/campaign"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/mango"

	"go.mongodb.org/mongo-driver/bson"
)

const CampaignCollection mango.Collection = "campaigns"

type campaignrepository struct {
	*MongoRepository
}

func Campaign() *campaignrepository {
	return &campaignrepository{
		New(CampaignCollection),
	}
}

// Getters & Setters
func (r *campaignrepository) GetDb() *mango.Mongo {
	return r.getDb()
}

func (r *campaignrepository) SetDb(db *mango.Mongo) *campaignrepository {
	r.setDb(db)
	return r
}

// Methods
func (r *campaignrepository) Load() ([]campaign.Campaign, error) {
	return r.fromCollection[campaign.Campaign](nil)
}

func (r *campaignrepository) FromCollectionByID(campaignID id.ID) (*campaign.Campaign, error) {
	campaigns, err := r.fromCollection[campaign.Campaign](bson.M{"id": campaignID})
	if err != nil {
		return nil, err
	}

	if len(campaigns) == 0 {
		return nil, logs.Error("error.not_found.id", map[string]any{
			"Subject": "Campaign",
			"ID":      campaignID,
		})
	}

	if len(campaigns) > 1 {
		return nil, logs.Error("error.unexpected:value", map[string]any{
			"Subject":  "Campaign",
			"Expected": 1,
			"Value":    len(campaigns),
		})
	}

	return &campaigns[0], nil
}

func (r *campaignrepository) Register(save *campaign.Campaign) error {
	return r.toCollection(save)
}
