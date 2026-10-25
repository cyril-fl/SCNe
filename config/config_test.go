package config

import (
	"os"
	"path/filepath"
	"scne/app/shared/services/i19n"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestPrivateConfigIsValid(t *testing.T) {
	var cfg config
	if err := yaml.Unmarshal(privateFile, &cfg); err != nil {
		t.Fatalf("Unmarshal(privateFile) error = %v", err)
	}
	if err := cfg.I18n.Validate(); err != nil {
		t.Fatalf("I18n.Validate() error = %v", err)
	}
	if cfg.Name == "" || cfg.Commands.Play.Use == "" {
		t.Fatalf("Name = %q, Commands.Play.Use = %q; want both set", cfg.Name, cfg.Commands.Play.Use)
	}
}

func TestMergeLocales(t *testing.T) {
	private := []i19n.Locale{{Code: "en", Name: "English"}, {Code: "fr", Name: "Français"}}
	public := []i19n.Locale{{Code: "zh", Name: "中文"}, {Code: "fr", Name: "French"}}

	got := mergeLocales(private, public)

	want := []i19n.Locale{{Code: "en", Name: "English"}, {Code: "fr", Name: "French"}, {Code: "zh", Name: "中文"}}
	if len(got) != len(want) {
		t.Fatalf("mergeLocales() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("mergeLocales()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
	if private[1].Name != "Français" {
		t.Fatalf("private locales modified: %v", private)
	}
}

func TestLoadPublicOverridesOnlyPublicKeys(t *testing.T) {
	cfg := config{Name: "scne"}
	cfg.I18n.Default = "fr"
	cfg.I18n.Locales = []i19n.Locale{{Code: "fr", Name: "Français"}}
	cfg.Theme.Accent = "#111111"
	cfg.Theme.Muted = "#222222"
	cfg.Documents.Files.Dice = "private/dice.yaml"

	loadPublic(&cfg, writeConfigFile(t, `
name: other
i18n:
  default: zh
  locales:
    - code: zh
      name: 中文
theme:
  accent: "#000000"
`))

	if cfg.Name != "scne" {
		t.Fatalf("Name = %q, want private value %q", cfg.Name, "scne")
	}
	if cfg.I18n.Default != "zh" || cfg.Theme.Accent != "#000000" {
		t.Fatalf("Default = %q, Accent = %q; want public values", cfg.I18n.Default, cfg.Theme.Accent)
	}
	if cfg.Theme.Muted != "#222222" || cfg.Documents.Files.Dice != "private/dice.yaml" {
		t.Fatalf("Muted = %q, Dice = %q; want private values", cfg.Theme.Muted, cfg.Documents.Files.Dice)
	}
	if len(cfg.I18n.Locales) != 2 || cfg.I18n.Locales[1].Code != "zh" {
		t.Fatalf("Locales = %v, want fr then zh", cfg.I18n.Locales)
	}
}

func TestLoadPublicPanicsOnInvalidYAML(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("loadPublic() did not panic on invalid YAML")
		}
	}()

	loadPublic(&config{}, writeConfigFile(t, "theme: [unclosed"))
}

// Helpers

func writeConfigFile(t *testing.T, contents string) string {
	t.Helper()

	file := filepath.Join(t.TempDir(), "public.yaml")
	if err := os.WriteFile(file, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}

	return file
}
