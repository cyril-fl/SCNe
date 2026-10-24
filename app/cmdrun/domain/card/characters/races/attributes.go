package races

import (
	"solopg/app/cmdrun/domain/card/attributes/stats"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

// - Race - //
// ID is stable across renames: characters reference their race by ID.
// Template races use their name as ID, since they cannot be renamed.
type Race struct {
	plumprimitive.Descriptible

	ID       string `bson:"id"`
	Playable bool
	Bonus    []stats.Modifier
}

type Template = struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Playable    bool             `json:"playable"`
	Bonus       []stats.Modifier `json:"bonus"`
}

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

func (r Race) GetID() string {
	return r.ID
}

func (r Race) IsPlayable() bool {
	return r.Playable
}

func (r Race) GetBonus() []stats.Modifier {
	return r.Bonus
}

func Assert(provided string) bool {
	return provided != "" && FindByID(provided) != nil
}
