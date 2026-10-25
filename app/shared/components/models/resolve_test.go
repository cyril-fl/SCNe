package models

import (
	"scne/app/shared/messages"
	"testing"
)

func TestResolveInitReturnsCompletedResolution(t *testing.T) {
	msg := Resolve().Init()()
	resolution, ok := msg.(messages.Resolved)
	if !ok || !resolution.Completed {
		t.Fatalf("unexpected resolution message: %#v", msg)
	}
}
