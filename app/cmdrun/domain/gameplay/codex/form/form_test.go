package codexform

import (
	"testing"

	"scne/app/shared/services/logs"
)

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

func TestValidateStatAcceptsIntegersOrEmpty(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })

	for _, value := range []string{"", "7", "-2"} {
		if err := validateStat(value); err != nil {
			t.Fatalf("validateStat(%q) = %v, want nil", value, err)
		}
	}
	for _, value := range []string{"abc", "1.5", " 3"} {
		if err := validateStat(value); err == nil {
			t.Fatalf("validateStat(%q) = nil, want an error", value)
		}
	}
}
