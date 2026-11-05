package cmdtemplate

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// store reads and writes the entries of a kind, by name.
type store interface {
	// Names returns every entry name.
	Names() ([]string, error)
	// Get returns the YAML of the entry name, or errNotFound.
	Get(name string) ([]byte, error)
	// Put replaces the entry name with data, or adds data when name is missing.
	Put(name string, data []byte) error
	// Delete removes the entry name, or returns errNotFound.
	Delete(name string) error
}

var errNotFound = errors.New("not found")

// - List file - //

// listStore keeps every entry in one YAML sequence file.
type listStore struct {
	path string
}

func (s listStore) Names() ([]string, error) {
	_, list, err := s.load()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(list.Content))
	for _, item := range list.Content {
		names = append(names, nameOf(item))
	}
	return names, nil
}

func (s listStore) Get(name string) ([]byte, error) {
	_, list, err := s.load()
	if err != nil {
		return nil, err
	}

	index := indexOf(list, name)
	if index < 0 {
		return nil, errNotFound
	}
	return encode(list.Content[index])
}

func (s listStore) Put(name string, data []byte) error {
	doc, list, err := s.load()
	if err != nil {
		return err
	}

	var item yaml.Node
	if err := yaml.Unmarshal(data, &item); err != nil {
		return err
	}
	if len(item.Content) == 0 {
		return errors.New("empty entry")
	}

	if index := indexOf(list, name); index >= 0 {
		// The comment above an entry, like a file header, belongs to the file, not to the entry.
		item.Content[0].HeadComment = list.Content[index].HeadComment
		list.Content[index] = item.Content[0]
	} else {
		list.Content = append(list.Content, item.Content[0])
	}
	return s.save(doc)
}

func (s listStore) Delete(name string) error {
	doc, list, err := s.load()
	if err != nil {
		return err
	}

	index := indexOf(list, name)
	if index < 0 {
		return errNotFound
	}
	list.Content = slices.Delete(list.Content, index, index+1)
	return s.save(doc)
}

// load returns the file document and its sequence; a missing file is an empty sequence.
func (s listStore) load() (*yaml.Node, *yaml.Node, error) {
	list := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	doc := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{list}}

	data, err := os.ReadFile(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		return doc, list, nil
	}
	if err != nil {
		return nil, nil, err
	}

	var loaded yaml.Node
	if err := yaml.Unmarshal(data, &loaded); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", s.path, err)
	}
	if len(loaded.Content) == 0 {
		return doc, list, nil
	}
	if loaded.Content[0].Kind != yaml.SequenceNode {
		return nil, nil, fmt.Errorf("%s: expected a list of entries", s.path)
	}
	return &loaded, loaded.Content[0], nil
}

func (s listStore) save(doc *yaml.Node) error {
	data, err := encode(doc)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

func indexOf(list *yaml.Node, name string) int {
	return slices.IndexFunc(list.Content, func(item *yaml.Node) bool { return nameOf(item) == name })
}

// nameOf returns the "name" value of a mapping node, or "".
func nameOf(node *yaml.Node) string {
	if node.Kind != yaml.MappingNode {
		return ""
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == "name" {
			return node.Content[i+1].Value
		}
	}
	return ""
}

// encode writes node as YAML indented by two spaces, like the data files.
func encode(node *yaml.Node) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	if err := encoder.Encode(node); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// - Folder - //

// folderStore keeps one entry per YAML file, in a folder and its subfolders.
type folderStore struct {
	dir string
}

func (s folderStore) Names() ([]string, error) {
	files, err := s.files()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	slices.Sort(names)
	return names, nil
}

func (s folderStore) Get(name string) ([]byte, error) {
	files, err := s.files()
	if err != nil {
		return nil, err
	}

	path, ok := files[name]
	if !ok {
		return nil, errNotFound
	}
	return os.ReadFile(path)
}

// Put writes over the file of name, or creates <dir>/<name>.yaml.
func (s folderStore) Put(name string, data []byte) error {
	files, err := s.files()
	if err != nil {
		return err
	}

	path, ok := files[name]
	if !ok {
		path = filepath.Join(s.dir, fileName(name))
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists", path)
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (s folderStore) Delete(name string) error {
	files, err := s.files()
	if err != nil {
		return err
	}

	path, ok := files[name]
	if !ok {
		return errNotFound
	}
	return os.Remove(path)
}

// files maps each entry name to its file; a missing folder has no entry.
func (s folderStore) files() (map[string]string, error) {
	files := map[string]string{}

	err := filepath.WalkDir(s.dir, func(path string, entry fs.DirEntry, err error) error {
		if errors.Is(err, fs.ErrNotExist) && path == s.dir {
			return fs.SkipAll
		}
		if err != nil || entry.IsDir() || !isYAML(path) {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var node yaml.Node
		if err := yaml.Unmarshal(data, &node); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if len(node.Content) > 0 {
			files[nameOf(node.Content[0])] = path
		}
		return nil
	})

	return files, err
}

func isYAML(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}

var unsafeChars = regexp.MustCompile(`[^a-z0-9._-]+`)

// fileName turns an entry name into a file name: "Black smith" gives "black_smith.yaml".
func fileName(name string) string {
	base := strings.Trim(unsafeChars.ReplaceAllString(strings.ToLower(name), "_"), "_.")
	if base == "" {
		base = "entry"
	}
	return base + ".yaml"
}
