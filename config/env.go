package config

import plumconfig "github.com/cyril-fl/plum/config"

var env = ""

func IsDev() bool {
	return plumconfig.EnvIs(env, "dev", "development", "develop", "")
}
