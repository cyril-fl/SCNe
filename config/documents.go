package config

type structureDocument struct {
	Folders structureFolders `yaml:"folders" json:"folders"`
	Files   structureFiles   `yaml:"files" json:"files"`
}

type structureFolders struct {
	Characters string `yaml:"characters" json:"characters"`
	Encounters string `yaml:"encounters" json:"encounters"`
	Hint       string `yaml:"hint" json:"hint"`
	Locations  string `yaml:"locations" json:"locations"`
	Oracle     string `yaml:"oracle" json:"oracle"`
}

type structureFiles struct {
	Classes         string `yaml:"classes" json:"classes"`
	Combat          string `yaml:"combat" json:"combat"`
	Dice            string `yaml:"dice" json:"dice"`
	ObjectsCategory string `yaml:"objects_category" json:"objects_category"`
	Priorities      string `yaml:"priorities" json:"priorities"`
	Ranks           string `yaml:"ranks" json:"ranks"`
	Races           string `yaml:"races" json:"races"`
	Rarity          string `yaml:"rarity" json:"rarity"`
	Slots           string `yaml:"slots" json:"slots"`
	Stats           string `yaml:"stats" json:"stats"`
	Variety         string `yaml:"variety" json:"variety"`
}
