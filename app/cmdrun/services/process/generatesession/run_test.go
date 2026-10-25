package generatesession

import (
	"testing"

	cmdruntui "scne/app/cmdrun/tui"
	"scne/app/shared/services/mango"
)

func TestProcessStoresDependencies(t *testing.T) {
	db := mango.New()
	ctx := &cmdruntui.Context{}
	p := Process(db, ctx)
	if p.db != db || p.ctx != ctx {
		t.Fatal("process did not retain its dependencies")
	}
}
