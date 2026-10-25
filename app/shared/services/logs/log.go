package logs

import (
	"fmt"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/style"
	"slices"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	plumprimitive "github.com/cyril-fl/plum/primitive"
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

var kindStyles = map[Kind]lipgloss.Style{
	SUCC: lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
	ERR:  lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
	INFO: lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
	WARN: lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
}

// -- Log -- //
type Log struct {
	plumprimitive.Datable

	ID      id.ID
	Type    Kind
	Author  string
	Message string
}

type Template struct {
	Type    Kind
	Author  string // Player ID or system
	Message string
}

func New(params Template) *Log {
	if params.Author == "" {
		params.Author = caches_author
	}

	return &Log{
		Datable: plumprimitive.Steady(),

		ID:      id.New(),
		Type:    params.Type,
		Author:  params.Author,
		Message: params.Message,
	}
}

// -- Methods -- //
func (l Log) String() string {
	kind := kindStyles[l.Type].Bold(true).Width(8).Render(strings.ToUpper(string(l.Type)))
	return fmt.Sprintf("%s %s%s %s", style.Muted().Render(l.GetFormattedUpdated(time.DateTime)), kind, style.Muted().Render(l.Author), l.Message)
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
	// Package init() (e.g. config) can log before logs.Init sets the registrable.
	if cache_registrable == nil {
		return
	}

	log := New(params)
	cache_registrable.Register(log)
}
