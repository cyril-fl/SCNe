package repository

import (
	"context"
	"scne/app/shared/services/logs"
	"scne/app/shared/services/mango"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const LoggingCollection mango.Collection = "log_system"

// Methods
type logsystemrepository struct {
	*MongoRepository
}

func LogSystem() *logsystemrepository {
	return &logsystemrepository{
		// Logs are written synchronously from the TUI loop: a down Mongo must not freeze it for 10s.
		New(LoggingCollection).setTimeout(time.Second),
	}
}

// Getters & Setters
func (r *logsystemrepository) GetDb() *mango.Mongo {
	return r.getDb()
}

func (r *logsystemrepository) SetDb(db *mango.Mongo) *logsystemrepository {
	r.setDb(db)
	return r
}

// Methods
func (r *logsystemrepository) Load(
	filter bson.M,
	opts ...FindOption,
) ([]logs.Log, error) {
	return r.fromCollection[logs.Log](filter, opts...)
}

func (r *logsystemrepository) Watch(ctx context.Context) (<-chan logs.Log, error) {
	return r.watchCollection[logs.Log](ctx)
}

func (r *logsystemrepository) Register(log *logs.Log) error {
	return r.toCollection(log)
}
