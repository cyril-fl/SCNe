package repository

import (
	"errors"
	"os"
	"scne/app/cmdrun/domain/campaign"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/mango"

	"go.mongodb.org/mongo-driver/bson"
)

const ArchivesCollection mango.Collection = "archives"

type archivesrepository struct {
	*MongoRepository
}

func Archives() *archivesrepository {
	return &archivesrepository{
		New(ArchivesCollection),
	}
}

// Getters & Setters
func (r *archivesrepository) GetDb() *mango.Mongo {
	return r.getDb()
}

func (r *archivesrepository) SetDb(db *mango.Mongo) *archivesrepository {
	r.setDb(db)
	return r
}

// Methods
func (r *archivesrepository) Load() ([]campaign.Archives, error) {
	return r.fromCollection[campaign.Archives](nil)
}

func (r *archivesrepository) LoadByCampaignID(campaignID id.ID) (*campaign.Archives, error) {
	archives, err := r.fromCollection[campaign.Archives](bson.M{"campaignId": campaignID})
	if err != nil {
		cwd, cwderr := os.Getwd()
		return nil, logs.Error("error.unexpected", map[string]any{
			"Path":  cwd,
			"Error": errors.Join(cwderr, err),
		})
	}

	/*
		NOTE If no archives found, returning (nil, nil) is perfectly fine.
		That's why return type is *campaign.Archives and not campaign.Archives.
		If not the whished behavior, the error should be handled appropriately in
		the calling function.
	*/
	if len(archives) == 0 {
		return nil, nil
	}

	if len(archives) > 1 {
		return nil, logs.Error("error.unexpected:value", map[string]any{
			"Subject":  "Archives",
			"Expected": 1,
			"Value":    len(archives),
		})
	}

	return &archives[0], nil
}

func (r *archivesrepository) Register(archives *campaign.Archives) error {
	return r.toCollection(archives)
}
