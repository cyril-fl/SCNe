package characters

import (
	"slices"
	"strings"
)

// ToggleCondition adds name to the conditions, or removes it when already there (case ignored).
// It reports whether name was added; a blank name changes nothing.
func (character *Character) ToggleCondition(name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}

	i := slices.IndexFunc(character.Conditions, func(c string) bool { return strings.EqualFold(c, name) })
	if i >= 0 {
		character.Conditions = slices.Delete(character.Conditions, i, i+1)
		return false
	}

	character.Conditions = append(character.Conditions, name)
	return true
}
