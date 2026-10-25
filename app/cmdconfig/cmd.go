package cmdconfig

import (
	"os"
	"os/exec"
	"scne/config"
	"strings"

	plumdebug "github.com/cyril-fl/plum/utils/debug"

	"github.com/spf13/cobra"
)

func RunCmdConfig(cobra *cobra.Command, args []string) {
	plumdebug.Json(config.GetConfig())
}

// RunECmdConfigEdit opens the public config file in the user editor, then validates it.
func RunECmdConfigEdit(cobra *cobra.Command, args []string) error {
	path, err := config.PublicFile()
	if err != nil {
		return err
	}

	if err := config.EnsurePublicFile(path); err != nil {
		return err
	}

	command := append(editor(), path)
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return config.Check()
}

// editor returns $VISUAL, then $EDITOR, then vi, split so "code --wait" works.
func editor() []string {
	for _, name := range []string{"VISUAL", "EDITOR"} {
		if fields := strings.Fields(os.Getenv(name)); len(fields) > 0 {
			return fields
		}
	}

	return []string{"vi"}
}
