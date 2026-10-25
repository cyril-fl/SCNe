package cmdconfig

import (
	"io"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	plumcliconfig "github.com/cyril-fl/plum/cli/config"
	plumcliutilsargs "github.com/cyril-fl/plum/cli/utils/args"
	"github.com/spf13/cobra"
)

func TestEditorPrefersVisualThenEditorThenVi(t *testing.T) {
	t.Setenv("VISUAL", "code --wait")
	t.Setenv("EDITOR", "nano")
	if got := editor(); !slices.Equal(got, []string{"code", "--wait"}) {
		t.Fatalf("editor() = %v, want VISUAL split", got)
	}

	t.Setenv("VISUAL", "")
	if got := editor(); !slices.Equal(got, []string{"nano"}) {
		t.Fatalf("editor() = %v, want EDITOR", got)
	}

	t.Setenv("EDITOR", "")
	if got := editor(); !slices.Equal(got, []string{"vi"}) {
		t.Fatalf("editor() = %v, want vi", got)
	}
}

func TestRunECmdConfigPrintsEachFormat(t *testing.T) {
	if got := runConfig(t); !strings.HasPrefix(got, "name: ") {
		t.Fatalf("default output = %q, want YAML", got)
	}
	if got := runConfig(t, "--json"); !strings.HasPrefix(got, "{\n  \"name\": ") {
		t.Fatalf("--json output = %q, want JSON with yaml key names", got)
	}
	if got := ansi.Strip(runConfig(t, "--pretty")); !strings.HasPrefix(got, "name: ") {
		t.Fatalf("--pretty output = %q, want YAML", got)
	}
	if got := ansi.Strip(runConfig(t, "--json", "--pretty")); !strings.HasPrefix(got, "{\n  \"name\": ") {
		t.Fatalf("--json --pretty output = %q, want JSON", got)
	}
}

// Helpers

func runConfig(t *testing.T, args ...string) string {
	t.Helper()

	cmd := &cobra.Command{}
	plumcliutilsargs.Makeflags(cmd, []plumcliconfig.Cmdflags{{Name: "yaml", Type: "bool"}, {Name: "json", Type: "bool"}, {Name: "pretty", Type: "bool"}})
	if err := cmd.ParseFlags(args); err != nil {
		t.Fatal(err)
	}

	stdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	t.Cleanup(func() { os.Stdout = stdout })

	if err := RunECmdConfig(cmd, nil); err != nil {
		t.Fatalf("RunECmdConfig() error = %v", err)
	}
	writer.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	return string(got)
}
