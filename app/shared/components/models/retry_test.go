package models

import (
	"testing"

	"charm.land/bubbles/v2/list"
)

func TestRerollModelRerollsAndTracksLimit(t *testing.T) {
	rolls := 0
	model := NewRerollModel(list.New(nil, list.NewDefaultDelegate(), 0, 0), func() (int, error) {
		rolls++
		return rolls, nil
	}, 2)

	if model.Attempt != 1 || model.Value != 1 || model.IsOutOfLimit() {
		t.Fatalf("unexpected initial model: attempt=%d value=%d", model.Attempt, model.Value)
	}
	model.Reroll()
	if model.Attempt != 2 || model.Value != 2 || !model.IsOutOfLimit() {
		t.Fatalf("unexpected rerolled model: attempt=%d value=%d", model.Attempt, model.Value)
	}
}
