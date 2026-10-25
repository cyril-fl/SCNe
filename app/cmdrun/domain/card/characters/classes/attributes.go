package classes

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"
	"solopg/app/cmdrun/domain/card/attributes/stats"
)

// - Class - //

// Class is a character class with its stat bonus and armor set.
type Class struct {
	name           string           `yaml:"name"`
	playable       bool             `yaml:"playable"`
	bonus          []stats.Modifier `yaml:"bonus"`
	equipementName string           `yaml:"armor_set"`
}

// Template configures a Class.
type Template = struct {
	Name           string           `yaml:"name"`
	Playable       bool             `yaml:"playable"`
	Bonus          []stats.Modifier `yaml:"bonus"`
	EquipementName string           `yaml:"armor_set"`
}

// New returns a Class.
func New(params Template) Class {
	return Class{
		name:           params.Name,
		playable:       params.Playable,
		bonus:          params.Bonus,
		equipementName: params.EquipementName,
	}
}

// GetName returns the class name.
func (c Class) GetName() string {
	return c.name
}

// IsPlayable reports whether the player can pick the class.
func (c Class) IsPlayable() bool {
	return c.playable
}

// GetBonus returns the stat bonus.
func (c Class) GetBonus() []stats.Modifier {
	return c.bonus
}

// GetEquipementName returns the armor set name.
func (c Class) GetEquipementName() string {
	return c.equipementName
}

// GetHash returns a SHA-256 of every class field.
func (c Class) GetHash() string {
	encoded, _ := json.Marshal(Template{
		Name:           c.GetName(),
		Playable:       c.IsPlayable(),
		Bonus:          c.GetBonus(),
		EquipementName: c.GetEquipementName(),
	})
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:])
}

// Assert reports whether provided is empty or a configured class name.
func Assert(provided string) bool {
	return provided == "" || slices.Contains(ListNames(), provided)
}
