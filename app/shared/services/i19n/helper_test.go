package i19n

import "testing"

func TestLocalizeWithoutConfiguredCacheReturnsPlaceholder(t *testing.T) {
	SetCache(nil)
	t.Cleanup(func() { SetCache(nil) })

	if got := Localize("missing.message"); got != "{{missing.message}}" {
		t.Fatalf("Localize() = %q, want placeholder", got)
	}
}

func TestMakeTemplateConvertsIntegerData(t *testing.T) {
	template := makeTemplate(nil, map[string]any{"count": 3, "value": true})
	if template["count"] != "3" || template["value"] != true {
		t.Fatalf("unexpected template: %#v", template)
	}
}
