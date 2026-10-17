package config

import (
	"os"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/services/yaml"
	"strings"

	"github.com/joho/godotenv"
)

var env = ""

type config struct {
	Name      string            `yaml:"name"`
	Verbose   bool              `yaml:"verbose"`
	I18n      i19n.Config       `yaml:"i18n"`
	Commands  commands          `yaml:"commands"`
	Documents structureDocument `yaml:"documents"`
}

type commands struct {
	Config  commandParams `yaml:"config"`
	Logs    commandParams `yaml:"logs"`
	Root    commandParams `yaml:"root"`
	Run     commandParams `yaml:"run"`
	Serve   commandParams `yaml:"serve"`
	Try     commandParams `yaml:"try"`
	Version commandParams `yaml:"version"`
}

type commandParams struct {
	Use     string `yaml:"use"`
	Short   string `yaml:"short"`
	Long    string `yaml:"long"`
	Example string `yaml:"example"`
	Args    []arg  `yaml:"args"`
}

type structureDocument struct {
	Folders structureFolders `yaml:"folders"`
	Files   structureFiles   `yaml:"files"`
}

type structureFolders struct {
	Characters string `yaml:"characters"`
	Encounters string `yaml:"encounters"`
	Hint       string `yaml:"hint"`
	Locations  string `yaml:"locations"`
	Oracle     string `yaml:"oracle"`
}

type structureFiles struct {
	Classes         string `yaml:"classes"`
	Dice            string `yaml:"dice"`
	ObjectsCategory string `yaml:"objects_category"`
	Priorities      string `yaml:"priorities"`
	Ranks           string `yaml:"ranks"`
	Races           string `yaml:"races"`
	Rarity          string `yaml:"rarity"`
	Slots           string `yaml:"slots"`
	Stats           string `yaml:"stats"`
	Variety         string `yaml:"variety"`
}

type arg struct {
	Name      string `yaml:"name"`
	Shorthand string `yaml:"shorthand"`
	Value     string `yaml:"value"`
	Usage     string `yaml:"usage"`
	Type      string `yaml:"type"`
}

func (c config) IsDev() bool {
	for _, validEnv := range []string{"dev", "development"} {
		if strings.EqualFold(env, validEnv) || env == "" {
			return true
		}
	}

	return false
}

type Arg = arg

var Current config

func load() (*config, error) {
	file := strings.TrimSpace(os.Getenv("CONFIG_FILE"))
	current, err := yaml.LoadFromFile[config](file)

	if err != nil {
		return nil, err
	}

	if err := current.I18n.Validate(); err != nil {
		return nil, logs.Error("error.invalid", map[string]any{
			"Subject": "config",
			"Error":   err,
		})
	}

	return current, nil
}

func init() {
	godotenv.Load()
	loaded, err := load()
	if err != nil {
		panic(err)
	}
	Current = *loaded
}
