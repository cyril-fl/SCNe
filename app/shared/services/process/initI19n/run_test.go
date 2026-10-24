package initI19n

import (
	"os"
	"path/filepath"
	"testing"

	"solopg/app/shared/services/i19n"
)

func testConfig(t *testing.T) i19n.Config {
	dir := t.TempDir()
	files := map[string]string{
		"en-US.yaml": "greeting: Hello\n",
		"fr-FR.yaml": "greeting: Bonjour\n",
		"de-DE.yaml": "greeting: Hallo\n",
		"notes.txt":  "ignored\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	t.Cleanup(func() { i19n.SetCache(i19n.New()) })

	return i19n.Config{
		Default: "en",
		Dir:     dir,
		Format:  i19n.FORMAT_YAML,
		Locales: []i19n.Locale{
			{Code: "en", ISO: "en-US", Name: "English", File: "en-US.yaml"},
			{Code: "fr", ISO: "fr-FR", Name: "Français", File: "fr-FR.yaml"},
		},
	}
}

func TestRunLoadsRequestedLocale(t *testing.T) {
	p := Process(testConfig(t), "fr")
	p.Run()

	if p.HasErr() {
		t.Fatalf("Run() error = %v", p.GetErr())
	}
	if len(p.files) != 2 {
		t.Fatalf("loaded files = %v, want only configured locales", p.files)
	}
	if got := i19n.Localize("greeting"); got != "Bonjour" {
		t.Fatalf("Localize() = %q, want Bonjour", got)
	}
	if got := i19n.Unlocalize("greeting"); got != "Hello" {
		t.Fatalf("Unlocalize() = %q, want Hello", got)
	}
}

func TestRunUsesDefaultLocaleWhenEmpty(t *testing.T) {
	p := Process(testConfig(t), "")
	p.Run()

	if got := i19n.Localize("greeting"); p.HasErr() || got != "Hello" {
		t.Fatalf("Localize() = %q, %v; want Hello, nil", got, p.GetErr())
	}
}

func TestRunFailsOnUnknownLocale(t *testing.T) {
	p := Process(testConfig(t), "es")
	p.Run()

	if !p.HasErr() {
		t.Fatal("Run() accepted an unknown locale")
	}
}

func TestRunFailsOnInvalidConfig(t *testing.T) {
	cfg := testConfig(t)
	cfg.Dir = filepath.Join(cfg.Dir, "missing")

	p := Process(cfg, "en")
	p.Run()

	if !p.HasErr() {
		t.Fatal("Run() accepted an invalid config")
	}
}
