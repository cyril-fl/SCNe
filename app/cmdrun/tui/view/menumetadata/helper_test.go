package menumetadata

import (
	"strings"
	"testing"

	"scne/app/cmdrun/domain/card"
	"scne/app/cmdrun/domain/card/attributes/stats"
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/card/characters/wallet"
	"scne/app/cmdrun/domain/card/locations"
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/shared/services/i19n"

	plumprimitive "github.com/cyril-fl/plum/primitive"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func TestSectionRendersLatestLocation(t *testing.T) {
	gameSession := &session.GameSession{State: &session.SessionState{}}
	section := NewSection(gameSession)

	gameSession.State.CurrentLocation = newLocation("forest")
	if view := section.GetView(); !strings.Contains(view, "forest") {
		t.Fatalf("expected view to contain %q, got %q", "forest", view)
	}

	gameSession.State.CurrentLocation = newLocation("desert")
	if view := section.GetView(); !strings.Contains(view, "desert") {
		t.Fatalf("expected view to contain %q, got %q", "desert", view)
	}
}

func TestSectionRendersCharacterSheet(t *testing.T) {
	useFrenchLocales(t)
	gameSession := &session.GameSession{State: &session.SessionState{
		CurrentLocation: newLocation("Taverne"),
		Player: &characters.Character{
			Card:   card.Card{Descriptible: plumprimitive.Describe(plumprimitive.TemplateDescriptible{Name: "Aragorn"})},
			Class:  "class.warrior",
			Race:   "race.human",
			Wallet: wallet.NewWallet(12, 5, 3),
			Stats:  stats.Stats{stats.List()[0]: 10},
		},
	}}

	view := NewSection(gameSession).GetView()
	t.Log("\n" + view)

	for _, want := range []string{
		"Nom      Aragorn",
		"Classe   Guerrier",
		"Race     Humain",
		"Lieu     Taverne",
		"Bourse   12 po · 5 pa · 3 pc",
		"Stats",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("view does not contain %q:\n%s", want, view)
		}
	}
}

func useFrenchLocales(t *testing.T) {
	t.Helper()
	bundle := i18n.NewBundle(language.French)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	for _, path := range []string{
		"locales/fr-FR.yaml",
		"locales/sidemenu/fr-FR.yaml",
		"locales/attributes/attributes-class/fr-FR.yaml",
		"locales/attributes/attributes-races/fr-FR.yaml",
		"locales/attributes/attributes-stats/fr-FR.yaml",
	} {
		if _, err := bundle.LoadMessageFile(path); err != nil {
			t.Fatalf("load %s: %v", path, err)
		}
	}
	translator := i19n.New()
	translator.SetLocale(i18n.NewLocalizer(bundle, "fr-FR"))
	i19n.SetCache(translator)
	t.Cleanup(func() { i19n.SetCache(nil) })
}

func newLocation(name string) *locations.Location {
	return &locations.Location{Card: card.Card{
		Descriptible: plumprimitive.Describe(plumprimitive.TemplateDescriptible{Name: name}),
	}}
}
