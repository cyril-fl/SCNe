package config

import plumcliconfig "github.com/cyril-fl/plum/cli/config"

type commands struct {
	Config     plumcliconfig.Cmdparams `yaml:"config"`
	ConfigEdit plumcliconfig.Cmdparams `yaml:"config_edit"`
	Logs       plumcliconfig.Cmdparams `yaml:"logs"`
	Root       plumcliconfig.Cmdparams `yaml:"root"`
	Play       plumcliconfig.Cmdparams `yaml:"play"`
	Serve      plumcliconfig.Cmdparams `yaml:"serve"`
	Try        plumcliconfig.Cmdparams `yaml:"try"`
	Version    plumcliconfig.Cmdparams `yaml:"version"`
}
