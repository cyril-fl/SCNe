package millstats

import (
	"testing"

	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters/classes"
	"scne/app/cmdrun/domain/card/characters/races"
	"scne/app/shared/services/logs"
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

func TestGetBonusCombinesClassAndRace(t *testing.T) {
	classList, raceList := classes.List(), races.List()
	if len(classList) == 0 || len(raceList) == 0 {
		t.Skip("no class or race in the loaded config")
	}
	class, race := classList[0], raceList[0]

	bonus := GetBonus(class.GetName(), race.GetName())

	if want := len(class.GetBonus()) + len(race.GetBonus()); len(bonus) != want {
		t.Fatalf("GetBonus returned %d modifiers, want %d", len(bonus), want)
	}
	if bonus := GetBonus("class.missing", "race.missing"); len(bonus) != 0 {
		t.Fatalf("GetBonus with unknown names = %#v, want empty", bonus)
	}
}
