package races

import (
	"solopg/app/cmdrun/domain/card/attributes/stats"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

// - Race - //

// Race is a character race with its stat bonus.
// ID is stable across renames: characters reference their race by ID.
// Template races use their name as ID, since they cannot be renamed.
type Race struct {
	plumprimitive.Descriptible

	ID       string `bson:"id"`
	Playable bool
	Bonus    []stats.Modifier
}

// Template configures a Race.
type Template = struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Playable    bool             `json:"playable"`
	Bonus       []stats.Modifier `json:"bonus"`
}

// New returns a Race.
func New(params Template) Race {
	return Race{
		Descriptible: plumprimitive.Describe(plumprimitive.TemplateDescriptible{
			Name:        params.Name,
			Description: params.Description,
		}),
		ID:       params.ID,
		Playable: params.Playable,
		Bonus:    params.Bonus,
	}
}

// GetID returns the race ID.
func (r Race) GetID() string {
	return r.ID
}

// IsPlayable reports whether the player can pick the race.
func (r Race) IsPlayable() bool {
	return r.Playable
}

// GetBonus returns the stat bonus.
func (r Race) GetBonus() []stats.Modifier {
	return r.Bonus
}

// Assert reports whether provided is a configured race ID.
func Assert(provided string) bool {
	return provided != "" && FindByID(provided) != nil
}
