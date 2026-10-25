package config

import (
	"os"
	"path/filepath"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/theme"
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

	err := loadPublic(&cfg, writeConfigFile(t, `
name: other
i18n:
  default: zh
  locales:
    - code: zh
      name: 中文
theme:
  accent: "#000000"
`))
	if err != nil {
		t.Fatalf("loadPublic() error = %v", err)
	}

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

func TestLoadPublicRejectsInvalidYAML(t *testing.T) {
	if err := loadPublic(&config{}, writeConfigFile(t, "theme: [unclosed")); err == nil {
		t.Fatal("loadPublic() error = nil, want unmarshal error")
	}
}

func TestLoadPublicIgnoresMissingFile(t *testing.T) {
	cfg := config{Name: "scne"}
	cfg.Theme.Accent = "#111111"

	if err := loadPublic(&cfg, filepath.Join(t.TempDir(), "missing.yaml")); err != nil {
		t.Fatalf("loadPublic() error = %v, want nil", err)
	}
	if cfg.Theme.Accent != "#111111" {
		t.Fatalf("Accent = %q, want unchanged %q", cfg.Theme.Accent, "#111111")
	}
}

func TestPublicTemplateChangesNothing(t *testing.T) {
	var public Public
	if err := yaml.Unmarshal(publicTemplate, &public); err != nil {
		t.Fatalf("Unmarshal(publicTemplate) error = %v", err)
	}
	if public.I18n.Default != "" || len(public.I18n.Locales) != 0 || public.Theme != (theme.Config{}) {
		t.Fatalf("publicTemplate sets values: %+v", public)
	}
}

func TestEnsurePublicFileCreatesTemplate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "scne", "config.yaml")

	if err := EnsurePublicFile(path); err != nil {
		t.Fatalf("EnsurePublicFile() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil || string(got) != string(publicTemplate) {
		t.Fatalf("file = %q, %v; want the template", got, err)
	}
}

func TestEnsurePublicFileKeepsExistingFile(t *testing.T) {
	path := writeConfigFile(t, "theme:\n  accent: \"#000000\"\n")

	if err := EnsurePublicFile(path); err != nil {
		t.Fatalf("EnsurePublicFile() error = %v", err)
	}

	got, _ := os.ReadFile(path)
	if string(got) != "theme:\n  accent: \"#000000\"\n" {
		t.Fatalf("file = %q, want unchanged", got)
	}
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
