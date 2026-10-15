package millstats

import (
	"testing"

	"solopg/app/cmdrun/domain/card/attributes/stats"
	"solopg/app/shared/services/logs"
)

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

func initTestLogs(t *testing.T) {
	t.Helper()
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })
}

func TestMakeModifierWithRandomFallbackUsesExplicitValue(t *testing.T) {
	modifier := makeModifierWithRandomFallback(withFallbackTemplate{
		stat:  stats.Stat("strength"),
		value: "7",
	})

	if modifier.Stat != stats.Stat("strength") || modifier.Value != 7 {
		t.Fatalf("got %#v, want strength: 7", modifier)
	}
}

func TestMakeModifiersFromOracleRejectsUnknownOracle(t *testing.T) {
	initTestLogs(t)
	modifiers, err := makeModifiersFromOracle("does-not-exist")
	if err == nil {
		t.Fatal("expected an error for an unknown oracle")
	}
	if modifiers != nil {
		t.Fatalf("modifiers = %#v, want nil", modifiers)
	}
}

func TestGenerateFromValuesResetsPreviousModifiers(t *testing.T) {
	initTestLogs(t)
	m := NewWithValues(map[string]string{"encounter": "does-not-exist"})
	m.modifiers = []stats.Modifier{{Stat: stats.Stat("strength"), Value: 9}}
	m.GenerateFromValues(FromValuesParams{})

	if len(m.modifiers) != 0 {
		t.Fatalf("modifiers = %#v, want empty after failed generation", m.modifiers)
	}
	if !m.HasErr() {
		t.Fatal("expected generation error")
	}
}
