package mango

import (
	"os"
	"testing"

	"solopg/config"
)

func TestReadEnvTrimsMongoValues(t *testing.T) {
	t.Setenv("MONGO_USERNAME", " user ")
	t.Setenv("MONGO_PASSWORD", " password ")
	t.Setenv("MONGO_PORT", " 27018 ")
	t.Setenv("MONGO_HOST", " db.example ")
	t.Setenv("MONGO_DBNAME", " campaign ")

	env := NewEnvironment()
	env.ReadEnv()

	if env.username != "user" || env.password != "password" || env.port != "27018" || env.host != "db.example" || env.dbname != "campaign" {
		t.Fatalf("ReadEnv() = %#v, want trimmed Mongo settings", env)
	}
}

func TestGetEnvURINormalizesCredentialsAndUsesDefaults(t *testing.T) {
	env := &MongoEnvironment{username: "user", password: "password"}

	if got, want := env.GetEnvURI(), "mongodb://user:password@localhost:27017"; got != want {
		t.Fatalf("GetEnvURI() = %q, want %q", got, want)
	}
	if env.host != HOST || env.port != PORT {
		t.Fatalf("GetEnvURI() did not apply defaults: host=%q port=%q", env.host, env.port)
	}
}

func TestGetEnvURIWithoutCredentials(t *testing.T) {
	env := &MongoEnvironment{host: "db.example", port: "27018"}

	if got, want := env.GetEnvURI(), "mongodb://db.example:27018"; got != want {
		t.Fatalf("GetEnvURI() = %q, want %q", got, want)
	}
}

func TestGetEnvDBNameUsesConfiguredNameAsFallback(t *testing.T) {
	original := config.Current.Name
	t.Cleanup(func() { config.Current.Name = original })
	config.Current.Name = "solopg-test"

	env := NewEnvironment()
	if got, want := env.GetEnvDBName(), "solopg-test"; got != want {
		t.Fatalf("GetEnvDBName() = %q, want %q", got, want)
	}
	if env.dbname != "solopg-test" {
		t.Fatalf("GetEnvDBName() did not cache fallback name: %q", env.dbname)
	}

	env.dbname = "custom"
	if got := env.GetEnvDBName(); got != "custom" {
		t.Fatalf("GetEnvDBName() = %q after explicit name, want custom", got)
	}
}

func lookupEnv(t *testing.T, key string) string {
	t.Helper()
	value, ok := os.LookupEnv(key)
	if !ok {
		t.Fatalf("environment variable %s was not set", key)
	}
	return value
}
