package resolvearchives

import (
	"testing"

	"scne/app/cmdrun/domain/campaign"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/logs"
)

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

func TestEnsureArchivesCreatesAndInitializesArchives(t *testing.T) {
	logs.Init(testLogRegistrar{})
	t.Cleanup(func() { logs.Init(nil) })
	p := Process(nil, id.New())

	p.ensureArchives()
	p.ensureArchivesInitialized()

	archives := p.GetResult()
	if archives == nil {
		t.Fatal("expected archives to be created")
	}
	if archives.CampaignID != p.campaignID {
		t.Fatal("archives campaign ID does not match process campaign ID")
	}
	if archives.Codex == nil || archives.Journal == nil || archives.Log == nil {
		t.Fatal("expected archives components to be initialized")
	}
}

func TestEnsureArchivesDoesNotReplaceLoadedArchives(t *testing.T) {
	p := Process(nil, id.New())
	p.cache.loadedArchives = &campaign.Archives{}

	p.ensureArchives()
	if p.GetResult().CampaignID != (id.ID{}) {
		t.Fatal("loaded archives were unexpectedly replaced")
	}
}
