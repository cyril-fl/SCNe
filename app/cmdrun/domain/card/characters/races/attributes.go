package races

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"
	"solopg/app/cmdrun/domain/card/attributes/stats"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

// - Race - //
type Race struct {
	plumprimitive.Descriptible

	Playable bool
	Bonus    []stats.Modifier
}

type Template = struct {
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
		Playable: params.Playable,
		Bonus:    params.Bonus,
	}
}

func (r Race) IsPlayable() bool {
	return r.Playable
}

func (r Race) GetBonus() []stats.Modifier {
	return r.Bonus
}

func (r Race) GetHash() string {
	encoded, _ := json.Marshal(Template{
		Name:        r.GetName(),
		Description: r.GetDescription(),
		Playable:    r.IsPlayable(),
		Bonus:       r.GetBonus(),
	})
	digest := sha256.Sum256(encoded)

	return hex.EncodeToString(digest[:])
}

func Assert(provided string) bool {
	return provided != "" || slices.Contains(ListNames(), provided)
}
