package id

import (
	"encoding/base64"

	"github.com/google/uuid"
)

// ID is a UUID identifier.
type ID uuid.UUID

// New returns a random ID.
func New() ID {
	return ID(uuid.New())
}

// IsZero reports whether id is the zero value.
func (id ID) IsZero() bool {
	return id == ID{}
}

// String returns the UUID text form.
func (id ID) String() string {
	return uuid.UUID(id).String()
}

// Base64 returns the standard base64 of the UUID bytes.
func (id ID) Base64() string {
	value := uuid.UUID(id)
	return base64.StdEncoding.EncodeToString(value[:])
}
