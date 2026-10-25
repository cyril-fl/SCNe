package mango

import (
	"net/url"
	"os"
	"scne/config"
	"strings"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// MongoEnvironment holds the MONGO_* connection settings.
type MongoEnvironment struct {
	username string
	password string
	port     string
	host     string
	dbname   string
}

/*
NOTE HOST and PORT, are use as default values, but they could be
overwritten per build environment with LDFLAGS if omit in .env file
*/
const (
	HOST = "localhost"
	PORT = "27017"
)

// NewEnvironment returns empty settings; call ReadEnv to fill them.
func NewEnvironment() *MongoEnvironment {
	return &MongoEnvironment{}
}

// ReadEnv reads the MONGO_* environment variables.
func (env *MongoEnvironment) ReadEnv() {
	env.username = strings.TrimSpace(os.Getenv("MONGO_USERNAME"))
	env.password = strings.TrimSpace(os.Getenv("MONGO_PASSWORD"))
	env.port = strings.TrimSpace(os.Getenv("MONGO_PORT"))
	env.host = strings.TrimSpace(os.Getenv("MONGO_HOST"))
	env.dbname = strings.TrimSpace(os.Getenv("MONGO_DBNAME"))
}

// GetEnvURI returns the connection URI; credentials are added only when both are set.
func (env *MongoEnvironment) GetEnvURI() string {
	uri := plumtransform.CleanJoin("@", env.getCredentials(), env.getHostAt())
	return plumtransform.CleanJoin("", "mongodb://", uri)
}

// GetEnvDBName returns MONGO_DBNAME, defaulting to the app name.
func (env *MongoEnvironment) GetEnvDBName() string {
	if env.dbname == "" {
		env.dbname = config.GetConfig().Name
	}

	return env.dbname
}

func (env *MongoEnvironment) getCredentials() string {
	if env.username == "" || env.password == "" {
		return ""
	}
	// Escapes "@", ":", "/" and "%" that would break the URI.
	return url.UserPassword(env.username, env.password).String()
}

func (env *MongoEnvironment) getHostAt() string {
	if env.host == "" {
		env.host = HOST
	}
	if env.port == "" {
		env.port = PORT
	}

	return plumtransform.CleanJoin(":", env.host, env.port)
}
