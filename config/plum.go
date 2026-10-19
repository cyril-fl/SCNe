package config

import "strings"

// TODO Move to plum/cli/config
type Cmdparams struct {
	Use     string   `yaml:"use"`
	Short   string   `yaml:"short"`
	Long    string   `yaml:"long"`
	Example string   `yaml:"example"`
	Args    []Cmdarg `yaml:"args"`
}

type Cmdarg struct {
	Name      string `yaml:"name"`
	Shorthand string `yaml:"shorthand"`
	Value     string `yaml:"value"`
	Usage     string `yaml:"usage"`
	Type      string `yaml:"type"`
}

// TODO Move in plum/env
func EnvIs(current string, valid ...string) bool {
	for _, env := range valid {
		if strings.EqualFold(current, env) {
			return true
		}
	}
	return false
}
