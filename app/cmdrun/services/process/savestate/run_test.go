package savestate

import (
	"testing"

	"solopg/app/cmdrun/services/game"
	"solopg/app/shared/services/mango"
)

func TestProcessStoresDependencies(t *testing.T) {
	db := mango.New()
	engine := &game.Engine{}
	p := Process(db, engine)

	if p.db != db || p.engine != engine {
		t.Fatal("process did not retain its database and engine dependencies")
	}
}
