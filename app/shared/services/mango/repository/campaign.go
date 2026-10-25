package repository

import (
	"scne/app/cmdrun/domain/campaign"
	"scne/app/shared/services/mango"
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

func (r *campaignrepository) Register(save *campaign.Campaign) error {
	return r.toCollection(save)
}
