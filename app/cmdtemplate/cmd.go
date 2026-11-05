package cmdtemplate

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"slices"

	"scne/app/shared/services/editor"
	"scne/app/shared/services/i19n"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/strictyaml"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// AddKindFlags adds one --<kind> flag per template type, exactly one being required.
func AddKindFlags(cmd *cobra.Command) {
	names := make([]string, 0, len(kinds()))
	for _, k := range kinds() {
		cmd.Flags().Bool(k.Name, false, k.Path)
		names = append(names, k.Name)
	}
	cmd.MarkFlagsOneRequired(names...)
	cmd.MarkFlagsMutuallyExclusive(names...)
}

// RunECmdAdd opens a commented template in the editor, then adds the entry once it is valid.
func RunECmdAdd(cmd *cobra.Command, args []string) error {
	k, err := selectedKind(cmd)
	if err != nil {
		return err
	}

	template, err := k.Template()
	if err != nil {
		return err
	}

	data, name, err := editOnce(template, func(data []byte) (string, error) {
		name, err := k.Validate(data)
		if err != nil {
			return "", err
		}
		return name, ensureNew(k, name)
	}, retryCommand(cmd, k, false))
	if err != nil || data == nil {
		return err
	}

	if err := k.store().Put(name, data); err != nil {
		return err
	}
	return printMessage("template.added", map[string]any{"Kind": k.Name, "Name": identifier(name), "Path": k.Path})
}

// RunECmdImport validates every entry of the given files, then adds them all, or none.
// --force replaces the entries that already exist.
func RunECmdImport(cmd *cobra.Command, args []string) error {
	k, err := selectedKind(cmd)
	if err != nil {
		return err
	}
	force, _ := cmd.Flags().GetBool("force")

	type entry struct {
		name string
		data []byte
	}
	entries := []entry{}
	errs := []error{}

	for _, path := range args {
		items, err := readEntries(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", path, err))
			continue
		}

		for _, data := range items {
			name, err := k.Validate(data)
			if err == nil && !force {
				err = ensureNew(k, name)
			}
			if err == nil && slices.ContainsFunc(entries, func(e entry) bool { return e.name == name }) {
				err = logs.Error("error.invalid:already_exists", map[string]any{"Subject": k.Name, "Value": identifier(name)})
			}
			if err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", path, err))
				continue
			}
			entries = append(entries, entry{name: name, data: data})
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	for _, e := range entries {
		if err := k.store().Put(e.name, e.data); err != nil {
			return err
		}
	}
	return printMessage("template.imported", map[string]any{"Kind": k.Name, "Count": len(entries), "Path": k.Path})
}

// RunECmdExport prints the YAML of an entry.
func RunECmdExport(cmd *cobra.Command, args []string) error {
	k, err := selectedKind(cmd)
	if err != nil {
		return err
	}

	data, err := get(k, args[0])
	if err != nil {
		return err
	}

	_, err = cmd.OutOrStdout().Write(data)
	return err
}

// RunECmdEdit opens an entry in the editor and saves it once it is valid; it needs --force.
func RunECmdEdit(cmd *cobra.Command, args []string) error {
	k, err := selectedKind(cmd)
	if err != nil {
		return err
	}
	name := args[0]

	data, err := get(k, name)
	if err != nil {
		return err
	}
	if err := ensureForced(cmd, name); err != nil {
		return err
	}

	edited, newName, err := editOnce(data, func(data []byte) (string, error) {
		newName, err := k.Validate(data)
		if err != nil || newName == name {
			return newName, err
		}
		return newName, ensureNew(k, newName)
	}, retryCommand(cmd, k, true))
	if err != nil || edited == nil {
		return err
	}

	if err := k.store().Put(name, edited); err != nil {
		return err
	}
	return printMessage("template.updated", map[string]any{"Kind": k.Name, "Name": identifier(newName), "Path": k.Path})
}

// RunECmdDelete removes an entry; it needs --force.
func RunECmdDelete(cmd *cobra.Command, args []string) error {
	k, err := selectedKind(cmd)
	if err != nil {
		return err
	}
	name := args[0]

	if _, err := get(k, name); err != nil {
		return err
	}
	if err := ensureForced(cmd, name); err != nil {
		return err
	}

	if err := k.store().Delete(name); err != nil {
		return err
	}
	return printMessage("template.deleted", map[string]any{"Kind": k.Name, "Name": identifier(name), "Path": k.Path})
}

// - Helpers - //

// identifier is an entry name shown as typed: i19n translates plain string values.
type identifier string

func selectedKind(cmd *cobra.Command) (kind, error) {
	for _, k := range kinds() {
		if selected, _ := cmd.Flags().GetBool(k.Name); selected {
			return k, nil
		}
	}
	return kind{}, errors.New("missing template type flag")
}

// retryCommand returns the import command that saves a fixed draft of k.
func retryCommand(cmd *cobra.Command, k kind, force bool) string {
	command := cmd.Root().Name() + " import --" + k.Name
	if force {
		command += " --force"
	}
	return command
}

// get returns the entry name, or a not found error naming it.
func get(k kind, name string) ([]byte, error) {
	data, err := k.store().Get(name)
	if errors.Is(err, errNotFound) {
		return nil, logs.Error("error.not_found.id", map[string]any{"Subject": k.Name, "ID": identifier(name)})
	}
	return data, err
}

// ensureNew fails when an entry named name exists.
func ensureNew(k kind, name string) error {
	names, err := k.store().Names()
	if err != nil {
		return err
	}
	if slices.Contains(names, name) {
		return logs.Error("template.error:exists", map[string]any{"Kind": k.Name, "Name": identifier(name)})
	}
	return nil
}

// ensureForced protects every existing entry: changing it needs --force.
func ensureForced(cmd *cobra.Command, name string) error {
	if force, _ := cmd.Flags().GetBool("force"); force {
		return nil
	}
	return logs.Error("template.error:force", map[string]any{"Name": identifier(name)})
}

// editOnce opens data in the editor once and returns the edited YAML once check accepts it.
// It returns nil data when the user closes the editor without a change, or leaves only
// comments. On an invalid file, the draft is kept and the error tells how to import it
// with retry, an import command without the path.
func editOnce(data []byte, check func([]byte) (string, error), retry string) ([]byte, string, error) {
	file, err := os.CreateTemp("", "scne-*.yaml")
	if err != nil {
		return nil, "", err
	}
	path := file.Name()
	file.Close()

	if err := os.WriteFile(path, data, 0o600); err != nil {
		os.Remove(path)
		return nil, "", err
	}
	if err := editor.Open(path); err != nil {
		os.Remove(path)
		return nil, "", err
	}

	edited, err := os.ReadFile(path)
	if err != nil {
		os.Remove(path)
		return nil, "", err
	}

	name, err := check(edited)
	if bytes.Equal(edited, data) || errors.Is(err, strictyaml.ErrEmpty) {
		os.Remove(path)
		lipgloss.Println(i19n.Localize("template.cancelled"))
		return nil, "", nil
	}
	if err != nil {
		return nil, "", logs.Error("template.error:draft", map[string]any{
			"Error":   identifier(err.Error()),
			"Path":    identifier(path),
			"Command": identifier(retry + " " + path),
		})
	}

	os.Remove(path)
	return edited, name, nil
}

// readEntries returns each entry of the YAML file at path: the items of a list, or the file itself.
func readEntries(path string) ([][]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	if len(doc.Content) == 0 {
		return nil, strictyaml.ErrEmpty
	}
	if doc.Content[0].Kind != yaml.SequenceNode {
		return [][]byte{data}, nil
	}

	entries := make([][]byte, 0, len(doc.Content[0].Content))
	for _, item := range doc.Content[0].Content {
		entry, err := encode(item)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func printMessage(id string, data map[string]any) error {
	lipgloss.Println(i19n.Localize(id, data))
	return nil
}
