package menuinventory

import (
	"testing"

	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters"
)

func TestApplyVitalsLosesAndRecovers(t *testing.T) {
	if !health.Validate() {
		t.Skip("health is not in the loaded stats config")
	}
	player := &characters.Character{Stats: stats.Stats{health: 10}}

	if _, err := applyVitals(player, map[string]string{"stat": "health", "operation": lose, "amount": "4"}); err != nil || player.Current(health) != 6 {
		t.Fatalf("lose: err = %v, health = %d; want nil, 6", err, player.Current(health))
	}
	if _, err := applyVitals(player, map[string]string{"stat": "health", "operation": recover, "amount": "9"}); err != nil || player.Current(health) != 10 {
		t.Fatalf("recover: err = %v, health = %d; want nil, 10 (capped)", err, player.Current(health))
	}
	if _, err := applyVitals(nil, map[string]string{}); err == nil {
		t.Fatal("applyVitals without player, want an error")
	}
}

func TestValidateAmount(t *testing.T) {
	if validateAmount("0") == nil || validateAmount("") == nil || validateAmount("3") != nil {
		t.Fatal("validateAmount must accept only amounts of at least 1")
	}
}
