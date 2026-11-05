package editor

import (
	"os"
	"os/exec"
	"strings"
)

// Open opens path in the user editor and waits for it to close.
func Open(path string) error {
	command := append(command(), path)
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

// command returns $VISUAL, then $EDITOR, then vi, split so "code --wait" works.
func command() []string {
	for _, name := range []string{"VISUAL", "EDITOR"} {
		if fields := strings.Fields(os.Getenv(name)); len(fields) > 0 {
			return fields
		}
	}

	return []string{"vi"}
}
