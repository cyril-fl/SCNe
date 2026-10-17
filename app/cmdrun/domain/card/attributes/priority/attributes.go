package priority

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

type Priority struct {
	name   string
	weight int
}

type Template struct {
	Name   string
	Weight int
}

func New(params Template) *Priority {
	return &Priority{
		name:   params.Name,
		weight: params.Weight,
	}
}

func (p Priority) GetHash() string {
	encoded := p.name + "-" + strconv.Itoa(p.weight)
	digest := sha256.Sum256([]byte(encoded))
	return hex.EncodeToString(digest[:])
}

func (p Priority) GetName() string {
	return p.name
}

func (p Priority) GetWeight() int {
	return p.weight
}
