package yaml

import (
	"os"
	"path/filepath"
	"solopg/app/shared/services/logs"

	"gopkg.in/yaml.v3"
)

// TODO Modifier pour pouvoir unmarshal du JSON, utiliser embeding et injection
func loadYAMLFromFile(fileAddress string, target any) error {
	data, err := os.ReadFile(fileAddress)
	if err != nil {
		return logs.RawError("read %s: %w", fileAddress, err)

	}

	if err := yaml.Unmarshal(data, target); err != nil {
		return logs.RawError("unmarshal %s: %w", fileAddress, err)
	}

	return nil
}

type wantedEntries struct {
	directories bool
	files       bool
}

func getFolderEntries(folderPath string, wanted wantedEntries) ([]string, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, logs.RawError("%q: %w", folderPath, err)
	}

	var names []string

	for _, entry := range entries {
		isDirectory := entry.IsDir()
		skipFile := !wanted.files && !isDirectory
		skipFolder := !wanted.directories && isDirectory

		if skipFile {
			continue
		}

		if skipFolder {
			continue
		}

		names = append(names, entry.Name())
	}

	return names, nil
}

func GetFolderFiles(folderPath string) ([]string, error) {
	return getFolderEntries(folderPath, wantedEntries{directories: false, files: true})
}

func GetFolderDirectories(folderPath string) ([]string, error) {
	return getFolderEntries(folderPath, wantedEntries{directories: true, files: false})
}

func GetFilesFromSource(sourcePath string, recursive bool) ([]string, error) {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return nil, logs.RawError("%q: %w", sourcePath, err)
	}

	if !info.IsDir() {
		return []string{sourcePath}, nil
	}

	return getFilesFromDirectoryRecursively(sourcePath, recursive)
}

func getFilesFromDirectoryRecursively(folderPath string, recursive bool) ([]string, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, logs.RawError("%q: %w", folderPath, err)
	}

	files := make([]string, 0, len(entries))

	for _, entry := range entries {
		path := filepath.Join(folderPath, entry.Name())

		if !entry.IsDir() {
			files = append(files, path)
			continue
		}

		if !recursive {
			continue
		}

		subFiles, err := getFilesFromDirectoryRecursively(path, true)
		if err != nil {
			return nil, logs.RawError("%q: %w", path, err)
		}

		files = append(files, subFiles...)
	}

	return files, nil
}

func LoadFromFile[T any](fileAddress string) (*T, error) {
	var params T
	if err := loadYAMLFromFile(fileAddress, &params); err != nil {
		return nil, logs.RawError("load %s: %w", fileAddress, err)
	}

	return &params, nil
}

func LoadListFromFile[T any](fileAddress string) ([]T, error) {
	var params []T
	if err := loadYAMLFromFile(fileAddress, &params); err != nil {
		return nil, logs.RawError("unmarshal %s: %w", fileAddress, err)
	}

	return params, nil
}
