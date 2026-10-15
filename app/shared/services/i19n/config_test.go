package i19n

import "testing"

func TestLocaleParseTag(t *testing.T) {
	tag, err := (&Locale{ISO: "fr-FR"}).ParseTag()
	if err != nil || tag.String() != "fr-FR" {
		t.Fatalf("ParseTag() = %v, %v; want fr-FR, nil", tag, err)
	}
	if _, err := (&Locale{ISO: ""}).ParseTag(); err == nil {
		t.Fatal("expected invalid locale error")
	}
}

func TestConfigValidationRejectsDuplicateLocales(t *testing.T) {
	cfg := Config{Default: "fr", Dir: ".", Format: FORMAT_YAML, Locales: []Locale{
		{Code: "fr", ISO: "fr-FR", Name: "Français", File: "fr-FR.yaml"},
		{Code: "fr", ISO: "fr-FR", Name: "Français", File: "fr-FR.yaml"},
	}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected duplicate locale error")
	}
}
