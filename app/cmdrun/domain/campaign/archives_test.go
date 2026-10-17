package campaign

import (
	"os"
	"testing"

	"solopg/app/cmdrun/domain/gameplay/codex"
	"solopg/app/cmdrun/types/id"
	"solopg/app/shared/services/logs"
)

type testLogRegistrar struct{}

func (testLogRegistrar) Register(*logs.Log) error { return nil }

func TestMain(m *testing.M) {
	logs.Init(testLogRegistrar{})
	code := m.Run()
	logs.Init(nil)
	os.Exit(code)
}

func TestNewArchives(t *testing.T) {
	campaignID := id.New()
	codexValue := &codex.Codex{}
	journal := NewJournal(nil)
	log := NewJournal(nil)

	archives := NewArchives(ArchivesTemplate{
		CampaignID: campaignID,
		Codex:      codexValue,
		Journal:    journal,
		Log:        log,
	})

	if archives.CampaignID != campaignID || archives.Codex != codexValue || archives.Journal != journal || archives.Log != log {
		t.Fatal("NewArchives did not copy its template")
	}
	if archives.GetcreatedAt().IsZero() || archives.GetupdatedAt().IsZero() {
		t.Fatal("NewArchives did not initialize dates")
	}
}

func TestArchivesFilter(t *testing.T) {
	campaignID := id.New()
	archives := NewArchives(ArchivesTemplate{CampaignID: campaignID})

	if got := archives.Filter()["campaignId"]; got != campaignID {
		t.Fatalf("Filter campaignId = %v, want %v", got, campaignID)
	}
}

func TestArchivesEnsureInitialized(t *testing.T) {
	archives := &Archives{}
	archives.EnsureInitialized()

	if archives.Codex == nil || archives.Journal == nil || archives.Log == nil {
		t.Fatal("EnsureInitialized did not create archive components")
	}
	if archives.Journal.Entries == nil || archives.Log.Entries == nil {
		t.Fatal("EnsureInitialized did not initialize journal entries")
	}
}
