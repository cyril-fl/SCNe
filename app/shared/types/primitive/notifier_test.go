package primitive

import (
	"testing"

	"charm.land/bubbles/v2/spinner"
)

func TestNotifierStartsInactive(t *testing.T) {
	notifier := NewNotifier()

	if notifier.IsNotifying() {
		t.Fatal("new notifier is active")
	}
	if notifier.GetNotification() != "" {
		t.Fatalf("new notifier notification = %q, want empty", notifier.GetNotification())
	}
	if notifier.GetSpinnerView() == "" {
		t.Fatal("new notifier has no spinner view")
	}
}

func TestNotifierHandlesNotificationAndIgnoresOverlap(t *testing.T) {
	notifier := NewNotifier()

	if cmd := notifier.HandleNotification(100, "notification:test"); cmd == nil {
		t.Fatal("HandleNotification returned nil for an inactive notifier")
	}
	if !notifier.IsNotifying() {
		t.Fatal("notifier is not active after HandleNotification")
	}
	if got := notifier.GetNotification(); got != "notification:test" {
		t.Fatalf("notification = %q, want notification:test", got)
	}

	if cmd := notifier.HandleNotification(100, "notification:other"); cmd != nil {
		t.Fatal("HandleNotification returned a command for an overlapping notification")
	}
	if got := notifier.GetNotification(); got != "notification:test" {
		t.Fatalf("overlapping notification replaced %q with %q", "notification:test", got)
	}
}

func TestNotifierHandlesSpinnerOnlyWhenActive(t *testing.T) {
	notifier := NewNotifier()

	if cmd := notifier.HandleSpinnerTick(nil); cmd != nil {
		t.Fatal("HandleSpinnerTick returned a command while inactive")
	}

	notifier.SetNotifying(true)
	if cmd := notifier.HandleSpinnerTick(spinner.TickMsg{}); cmd == nil {
		t.Fatal("HandleSpinnerTick returned nil while active")
	}
}
