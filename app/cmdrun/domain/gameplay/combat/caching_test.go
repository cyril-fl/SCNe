package combat

import (
	"testing"

	"solopg/app/cmdrun/domain/card/characters/classes"
	"solopg/app/cmdrun/domain/gameplay/oracle"
	"solopg/app/shared/services/logs"
)

func TestValidateRejectsInvalidActions(t *testing.T) {
	t.Cleanup(func() { logs.Init(nil) })
	logs.Init(combatTestLogRegistry{})

	tests := []struct {
		name   string
		config yamlConfig
	}{
		{name: "missing id", config: yamlConfig{Oracle: "combat_attack"}},
		{name: "unknown oracle", config: yamlConfig{ID: "mock", Oracle: "mock.oracle"}},
		{name: "unknown effect", config: yamlConfig{ID: "mock", Oracle: "combat_attack", Effects: map[string]map[string]int{"hit": {"poison": 1}}}},
		{name: "negative value", config: yamlConfig{ID: "mock", Oracle: "combat_attack", Effects: map[string]map[string]int{"hit": {effectDamage: -1}}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := validate(test.config); err == nil {
				t.Fatal("expected an error")
			}
		})
	}

	valid := yamlConfig{ID: "mock", Oracle: "combat_attack", Effects: map[string]map[string]int{"hit": {effectDamage: 1, effectGuard: 0}}}
	if err := validate(valid); err != nil {
		t.Fatalf("validate(valid) error = %v", err)
	}
}

func TestClassCombatReferencesExist(t *testing.T) {
	if _, err := ListWithError(); err != nil {
		t.Fatalf("ListWithError() error = %v", err)
	}

	for _, class := range classes.List() {
		for _, id := range class.GetActions() {
			if _, err := GetByID(id); err != nil {
				t.Errorf("class %s: action %s: %v", class.GetName(), id, err)
			}
		}
		if class.GetIntent() == "" {
			continue
		}
		if _, err := oracle.GetByName(class.GetIntent()); err != nil {
			t.Errorf("class %s: intent %s: %v", class.GetName(), class.GetIntent(), err)
		}
	}
}
