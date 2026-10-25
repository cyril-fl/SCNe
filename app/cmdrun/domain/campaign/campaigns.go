package campaign

import (
	"fmt"

	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/card/characters/races"
	"scne/app/cmdrun/domain/card/locations"

	"scne/app/cmdrun/types/id"
	"scne/app/cmdrun/types/timer"
	"scne/app/shared/services/i19n"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	plumtransform "github.com/cyril-fl/plum/utils/transform"
	"go.mongodb.org/mongo-driver/bson"
)

// Campaign is a saved game: the player, the current location and the time played.
type Campaign struct {
	plumprimitive.Datable

	ID              id.ID                 `bson:"Id" json:"Id"`
	Player          *characters.Character `bson:"character" json:"character"`
	CurrentLocation *locations.Location   `bson:"location" json:"location"`
	TimeHistory     []timer.Record        `bson:"timeHistory" json:"timeHistory"`
}

// Template configures a Campaign.
type Template struct {
	Player          *characters.Character
	CurrentLocation *locations.Location
}

// New returns a dated Campaign with a new ID.
func New(params Template) *Campaign {
	return &Campaign{
		Datable: plumprimitive.Dated(),

		ID:              id.New(),
		Player:          params.Player,
		CurrentLocation: params.CurrentLocation,
	}
}

// Title returns "name | race | class" of the player.
func (c *Campaign) Title() string {
	return fmt.Sprintf("%s | %s | %s",
		c.Player.GetName(),
		plumtransform.Capitalize(i19n.Localize(races.NameOf(c.Player.Race))),
		plumtransform.Capitalize(i19n.Localize(c.Player.Class)),
	)
}

// Description returns the localized location and last update.
func (c *Campaign) Description() string {
	return i19n.Localize("campaign.description", map[string]any{
		"Location": c.CurrentLocation.GetName(),
		"Updated":  c.GetFormattedUpdated("2006-01-02 15:04:05"),
	})
}

// Filter returns the MongoDB filter matching the campaign.
func (c *Campaign) Filter() bson.M {
	return bson.M{
		"Id": c.ID,
	}
}
