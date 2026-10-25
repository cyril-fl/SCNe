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

func TestRunReportsCampaignErrorWithoutSession(t *testing.T) {
	p := Process(nil, &cmdruntui.Context{})
	p.Run()

	if !p.HasErr() {
		t.Fatal("an empty onboarding should fail")
	}
	if p.GetResult() != nil {
		t.Fatal("no session should be built when the campaign fails")
	}
}
