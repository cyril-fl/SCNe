package priority

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

// Priority is a named quest priority with a weight.
type Priority struct {
	name   string
	weight int
}

// Template configures a Priority.
type Template struct {
	Name   string
	Weight int
}

// New returns a Priority.
func New(params Template) *Priority {
	return &Priority{
		name:   params.Name,
		weight: params.Weight,
	}
}

// GetHash returns a SHA-256 of the name and weight.
func (p Priority) GetHash() string {
	encoded := p.name + "-" + strconv.Itoa(p.weight)
	digest := sha256.Sum256([]byte(encoded))
	return hex.EncodeToString(digest[:])
}

// GetName returns the priority name.
func (p Priority) GetName() string {
	return p.name
}

// GetWeight returns the priority weight.
func (p Priority) GetWeight() int {
	return p.weight
}
