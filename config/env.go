package config

var env = ""

func IsDev() bool {
	return EnvIs(env, "dev", "development", "develop", "")
}
