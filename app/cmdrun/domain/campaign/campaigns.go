package campaign

import (
	"fmt"

	"solopg/app/cmdrun/domain/card/characters"
	"solopg/app/cmdrun/domain/card/locations"

	"solopg/app/cmdrun/types/id"
	"solopg/app/cmdrun/types/timer"
	"solopg/app/shared/services/i19n"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plumtransform "github.com/cyril-fl/plum/utils/transform"
	"go.mongodb.org/mongo-driver/bson"
)

// TEST
type Campaign struct {
	plumprimitive.Datable

	ID              id.ID                 `bson:"Id" json:"Id"`
	Player          *characters.Character `bson:"character" json:"character"`
	CurrentLocation *locations.Location   `bson:"location" json:"location"`
	TimeHistory     []timer.Record        `bson:"timeHistory" json:"timeHistory"`
}

type Template struct {
	Player          *characters.Character
	CurrentLocation *locations.Location
}

func New(params Template) *Campaign {
	return &Campaign{
		Datable: plumprimitive.Dated(),

		ID:              id.New(),
		Player:          params.Player,
		CurrentLocation: params.CurrentLocation,
	}
}

func (c *Campaign) Title() string {
	return fmt.Sprintf("%s | %s | %s",
		c.Player.GetName(),
		plumtransform.Capitalize(i19n.Localize(c.Player.Race)),
		plumtransform.Capitalize(i19n.Localize(c.Player.Class)),
	)
}

func (c *Campaign) Description() string {
	return i19n.Localize("campaign.description", map[string]any{
		"Location": c.CurrentLocation.GetName(),
		"Updated":  c.GetFormattedUpdated("2006-01-02 15:04:05"),
	})
}

func (c *Campaign) Filter() bson.M {
	return bson.M{
		"Id": c.ID,
	}
}
