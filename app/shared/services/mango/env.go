package mango

import (
	"os"
	"solopg/config"
	"strings"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

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

func NewEnvironment() *MongoEnvironment {
	return &MongoEnvironment{}
}

func (env *MongoEnvironment) ReadEnv() {
	env.username = strings.TrimSpace(os.Getenv("MONGO_USERNAME"))
	env.password = strings.TrimSpace(os.Getenv("MONGO_PASSWORD"))
	env.port = strings.TrimSpace(os.Getenv("MONGO_PORT"))
	env.host = strings.TrimSpace(os.Getenv("MONGO_HOST"))
	env.dbname = strings.TrimSpace(os.Getenv("MONGO_DBNAME"))
}

func (env *MongoEnvironment) GetEnvURI() string {
	uri := plumtransform.CleanJoin("@", env.getCredentials(), env.getHostAt())
	return plumtransform.CleanJoin("", "mongodb://", uri)
}

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
	return plumtransform.CleanJoin(":", env.username, env.password)
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
