package cmdtemplate

import (
	"embed"
	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/card/characters/classes"
	"scne/app/cmdrun/domain/card/characters/races"
	"scne/app/cmdrun/domain/card/locations"
	"scne/app/cmdrun/domain/gameplay/hint"
	"scne/app/cmdrun/domain/gameplay/oracle"
	"scne/config"
)

//go:embed templates/*.yaml
var templates embed.FS

// kind is a template type: one flag, one config path and its validation.
type kind struct {
	Name string
	// Path is a list file when List is set, else a folder of one file per entry.
	Path     string
	List     bool
	Validate func(data []byte) (string, error)
}

// kinds lists the editable template types, in flag order.
func kinds() []kind {
	documents := config.GetConfig().Documents

	return []kind{
		{Name: "class", Path: documents.Files.Classes, List: true, Validate: classes.Validate},
		{Name: "race", Path: documents.Files.Races, List: true, Validate: races.Validate},
		{Name: "character", Path: documents.Folders.Characters, Validate: characters.Validate},
		{Name: "location", Path: documents.Folders.Locations, Validate: locations.Validate},
		{Name: "hint", Path: documents.Folders.Hint, Validate: hint.Validate},
		{Name: "oracle", Path: documents.Folders.Oracle, Validate: oracle.Validate},
		{Name: "encounter", Path: documents.Folders.Encounters, Validate: oracle.Validate},
	}
}

// Template returns the commented YAML opened by `add`.
func (k kind) Template() ([]byte, error) {
	return templates.ReadFile("templates/" + k.Name + ".yaml")
}

// store returns the storage of k's entries.
func (k kind) store() store {
	if k.List {
		return listStore{path: k.Path}
	}
	return folderStore{dir: k.Path}
}
