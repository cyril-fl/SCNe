package generatesession

import (
	"testing"

	cmdruntui "solopg/app/cmdrun/tui"
	"solopg/app/shared/services/mango"
)

func TestProcessStoresDependencies(t *testing.T) {
	db := mango.New()
	ctx := &cmdruntui.Context{}
	p := Process(db, ctx)
	if p.db != db || p.ctx != ctx {
		t.Fatal("process did not retain its dependencies")
	}
}
