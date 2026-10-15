package models

import (
	"testing"

	sharedtui "solopg/app/shared/tui"
)

func TestResolveInitReturnsCompletedResolution(t *testing.T) {
	msg := Resolve().Init()()
	resolution, ok := msg.(sharedtui.ResolutionMsg)
	if !ok || !resolution.Completed {
		t.Fatalf("unexpected resolution message: %#v", msg)
	}
}
