package savestate

import (
	"testing"

	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/shared/services/mango"
)

func TestProcessStoresDependencies(t *testing.T) {
	db := mango.New()
	gameSession := &session.GameSession{}
	p := Process(db, gameSession)

	if p.db != db || p.session != gameSession {
		t.Fatal("process did not retain its database and session dependencies")
	}
}
