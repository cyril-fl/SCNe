package generatebuildinfo

import (
	"strings"
	"testing"
)

func TestProcessUsesProductionBuildInfo(t *testing.T) {
	oldVersion, oldCommit, oldDate := Version, Commit, Date
	Version, Commit, Date = "1.2.3", "abc123", "2026-10-15T00:00:00Z"
	t.Cleanup(func() { Version, Commit, Date = oldVersion, oldCommit, oldDate })

	p := Process()
	p.makeProdEnv()
	got := p.GetResult()
	if got.String() == "" {
		t.Fatal("expected non-empty build info")
	}
	for _, value := range []string{Version, Commit, Date} {
		if !strings.Contains(got.String(), value) {
			t.Fatalf("build info %q does not contain %q", got.String(), value)
		}
	}
}
