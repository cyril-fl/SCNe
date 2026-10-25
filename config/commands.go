package config

import plumcliconfig "github.com/cyril-fl/plum/cli/config"

type commands struct {
	Config     plumcliconfig.Cmdparams `yaml:"config" json:"config"`
	ConfigEdit plumcliconfig.Cmdparams `yaml:"config_edit" json:"config_edit"`
	Logs       plumcliconfig.Cmdparams `yaml:"logs" json:"logs"`
	Root       plumcliconfig.Cmdparams `yaml:"root" json:"root"`
	Play       plumcliconfig.Cmdparams `yaml:"play" json:"play"`
	Serve      plumcliconfig.Cmdparams `yaml:"serve" json:"serve"`
	Try        plumcliconfig.Cmdparams `yaml:"try" json:"try"`
	Version    plumcliconfig.Cmdparams `yaml:"version" json:"version"`
}
