package campaign

import (
	"os"
	"testing"

	"scne/app/cmdrun/domain/card/attributes/rarity"
	"scne/app/cmdrun/domain/card/attributes/variety"
	"scne/app/cmdrun/domain/card/locations"
	"scne/app/cmdrun/domain/card/quest"
	"scne/app/cmdrun/domain/gameplay/codex"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/services/logs"

	"go.mongodb.org/mongo-driver/bson"
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

	archives := NewArchives(ArchivesTemplate{
		CampaignID: campaignID,
		Codex:      codexValue,
		Journal:    journal,
	})

	if archives.CampaignID != campaignID || archives.Codex != codexValue || archives.Journal != journal {
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

	if archives.Codex == nil || archives.Journal == nil {
		t.Fatal("EnsureInitialized did not create archive components")
	}
	if archives.Journal.Entries == nil {
		t.Fatal("EnsureInitialized did not initialize journal entries")
	}
}

func TestArchivesBSONRoundTripPreservesCodexEntries(t *testing.T) {
	location, err := locations.New(locations.Template{
		Name:        "Saved location",
		Description: "A location that must survive a save",
		Rarity:      rarity.Default(),
		Variety:     variety.MakeDefault("location_card"),
	})
	if err != nil {
		t.Fatalf("locations.New() error = %v", err)
	}

	newQuest, err := quest.New(quest.Template{
		Name:        "Saved quest",
		Description: "A quest that must survive a save",
	})
	if err != nil {
		t.Fatalf("quest.New() error = %v", err)
	}

	archives := NewArchives(ArchivesTemplate{
		CampaignID: id.New(),
		Codex:      codex.New(),
		Journal:    NewJournal([]Entry{}),
	})
	archives.Codex.Locations.Add(location)
	archives.Codex.Quest.Add(newQuest)

	encoded, err := bson.Marshal(archives)
	if err != nil {
		t.Fatalf("bson.Marshal() error = %v", err)
	}

	var decoded Archives
	if err := bson.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("bson.Unmarshal() error = %v", err)
	}
	decoded.EnsureInitialized()

	if got := len(decoded.Codex.Locations.FormEntries()); got != 1 {
		t.Fatalf("decoded locations = %d, want 1", got)
	}
	if got := len(decoded.Codex.Quest.FormEntries()); got != 1 {
		t.Fatalf("decoded quests = %d, want 1", got)
	}
}
