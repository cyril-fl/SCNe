package codextable

import (
	"testing"

	"scne/app/cmdrun/domain/card/locations"
)

func TestLocationEffectsRoundTrip(t *testing.T) {
	values := map[string]string{
		"effect.0.name":        "Fresh air",
		"effect.0.description": "Increases stamina by 1",
		"effect.0.stamina":     "1",
		"effect.1.name":        "",
		"effect.1.description": "",
	}

	effects := effectsFromValues(values)
	if len(effects) != 1 || effects[0].GetName() != "Fresh air" || effects[0].Modifier.Stat != "stamina" || effects[0].Modifier.Value != 1 {
		t.Fatalf("effectsFromValues() = %+v, want one stamina +1 effect", effects)
	}

	table := NewLocations()
	table.Add(&locations.Location{Effects: effects})
	got := table.FormEntries()[0].Values
	for _, key := range []string{"effect.0.name", "effect.0.description", "effect.0.stamina"} {
		if got[key] != values[key] {
			t.Fatalf("FormEntries()[%q] = %q, want %q", key, got[key], values[key])
		}
	}
	if _, ok := got["effect.1.name"]; ok {
		t.Fatal("FormEntries() kept the empty effect group")
	}
}
