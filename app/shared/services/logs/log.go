package logs

import (
	"fmt"
	"slices"
	"solopg/app/cmdrun/types/id"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/types/primitive"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// -- Kind -- //
type Kind string

const (
	SUCC Kind = "success"
	ERR  Kind = "error"
	INFO Kind = "info"
	WARN Kind = "warn"
)

func AssertKind(kind Kind) bool {
	kinds := []Kind{
		SUCC,
		ERR,
		INFO,
		WARN,
	}
	return slices.Contains(kinds, kind)
}

// -- Log -- //
type Log struct {
	primitive.Datable

	ID      id.ID
	Type    Kind
	Author  string
	Message string
}

type Template struct {
	Type    Kind
	Author  string // Player Id O system
	Message string
}

func New(params Template) *Log {
	if params.Author == "" {
		params.Author = caches_author
	}

	return &Log{
		Datable: primitive.Steady(),

		ID:      id.New(),
		Type:    params.Type,
		Author:  params.Author,
		Message: params.Message,
	}
}

// -- Methods -- //
func (l Log) String() string {
	return fmt.Sprintf("%s - [%s] - (%s) : \"%s\"", l.GetFormattedUpdated(time.RFC3339), strings.ToUpper(string(l.Type)), l.Author, l.Message)
}

func (l *Log) Touch() {
	// Logs are immutable and must keep the previous SetUpdatedAt behavior.
}

func (l *Log) Filter() bson.M {
	return bson.M{
		"id": l.ID,
	}
}

// -- Helpers -- //
// SystemLog creates a new log entry, registers the unlocalized message, and returns the localized one.
func SystemLog(id string, data ...map[string]any) string {
	registerFromTemplate(Template{
		Type:    ERR,
		Message: i19n.Unlocalize(id, data...),
	})

	return i19n.Localize(id, data...)
}

func registerFromTemplate(params Template) {
	log := New(params)
	cache_registrable.Register(log)
}
