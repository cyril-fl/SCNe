package config

type commands struct {
	Config  Cmdparams `yaml:"config"`
	Logs    Cmdparams `yaml:"logs"`
	Root    Cmdparams `yaml:"root"`
	Run     Cmdparams `yaml:"run"`
	Serve   Cmdparams `yaml:"serve"`
	Try     Cmdparams `yaml:"try"`
	Version Cmdparams `yaml:"version"`
}
