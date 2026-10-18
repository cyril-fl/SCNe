package id

import (
	"encoding/base64"

	"github.com/google/uuid"
)

// TODO rechercher où il n'est pas utilisé et uniformiser le comportements
type ID uuid.UUID

func New() ID {
	return ID(uuid.New())
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

func (id ID) Base64() string {
	value := uuid.UUID(id)
	return base64.StdEncoding.EncodeToString(value[:])
}
