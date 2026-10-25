package cmdrunrunsession

import (
	"testing"

	cmdruntui "scne/app/cmdrun/tui"
	"scne/app/shared/services/logs"
)

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

func TestMakeOnboardingSteps(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })
	steps := makeOnboardingSteps()
	if len(steps) != 5 {
		t.Fatalf("got %d onboarding steps, want 5", len(steps))
	}

	var ctx cmdruntui.Context
	if err := steps[0].Resolve(&ctx, "Aria"); err != nil {
		t.Fatalf("resolving name returned error: %v", err)
	}
	if ctx.SelectedName != "Aria" {
		t.Fatalf("selected name = %q, want Aria", ctx.SelectedName)
	}
	if err := steps[0].Resolve(&ctx, 42); err == nil {
		t.Fatal("expected an error for a non-string name")
	}
}
