package config

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
