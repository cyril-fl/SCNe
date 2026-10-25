package menumetadata

import (
	"strings"
	"testing"

	"scne/app/cmdrun/domain/card"
	"scne/app/cmdrun/domain/card/locations"
	"scne/app/cmdrun/domain/gameplay/session"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

func TestSectionRendersLatestLocation(t *testing.T) {
	gameSession := &session.GameSession{State: &session.SessionState{}}
	section := NewSection(gameSession)

	gameSession.State.CurrentLocation = newLocation("forest")
	if view := section.GetView(); !strings.Contains(view, "forest") {
		t.Fatalf("expected view to contain %q, got %q", "forest", view)
	}

	gameSession.State.CurrentLocation = newLocation("desert")
	if view := section.GetView(); !strings.Contains(view, "desert") {
		t.Fatalf("expected view to contain %q, got %q", "desert", view)
	}
}

func newLocation(name string) *locations.Location {
	return &locations.Location{Card: card.Card{
		Descriptible: plumprimitive.Describe(plumprimitive.TemplateDescriptible{Name: name}),
	}}
}
