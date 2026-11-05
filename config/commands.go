package config

import plumcliconfig "github.com/cyril-fl/plum/cli/config"

type commands struct {
	Add        plumcliconfig.Cmdparams `yaml:"add" json:"add"`
	Delete     plumcliconfig.Cmdparams `yaml:"delete" json:"delete"`
	Edit       plumcliconfig.Cmdparams `yaml:"edit" json:"edit"`
	Export     plumcliconfig.Cmdparams `yaml:"export" json:"export"`
	Import     plumcliconfig.Cmdparams `yaml:"import" json:"import"`
	Config     plumcliconfig.Cmdparams `yaml:"config" json:"config"`
	ConfigEdit plumcliconfig.Cmdparams `yaml:"config_edit" json:"config_edit"`
	Logs       plumcliconfig.Cmdparams `yaml:"logs" json:"logs"`
	Root       plumcliconfig.Cmdparams `yaml:"root" json:"root"`
	Play       plumcliconfig.Cmdparams `yaml:"play" json:"play"`
	Serve      plumcliconfig.Cmdparams `yaml:"serve" json:"serve"`
	Try        plumcliconfig.Cmdparams `yaml:"try" json:"try"`
	Version    plumcliconfig.Cmdparams `yaml:"version" json:"version"`
}
