package repository

import (
	"context"
	"time"

	"scne/app/shared/services/logs"
	"scne/app/shared/services/mango"

	"go.mongodb.org/mongo-driver/bson"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --- Repository --
// Document
type Document interface {
	Filter() bson.M
	Touch()
}

type MongoRepository struct {
	db         *mango.Mongo
	collection mango.Collection
	timeout    time.Duration
}

func New(collection mango.Collection) *MongoRepository {
	return &MongoRepository{
		collection: collection,
	}
}

// Getters & Setters
func (r *MongoRepository) getDb() *mango.Mongo {
	return r.db
}

func (r *MongoRepository) setDb(db *mango.Mongo) *MongoRepository {
	r.db = db
	return r
}

func (r *MongoRepository) getTimeout() time.Duration {
	if r.timeout > 0 {
		return r.timeout
	}
	return 10 * time.Second
}

func (r *MongoRepository) setTimeout(timeout time.Duration) *MongoRepository {
	r.timeout = timeout
	return r
}

func (r *MongoRepository) getCollection() mango.Collection {
	return r.collection
}

// Methods
func (r *MongoRepository) fromCollection[T any](
	filter bson.M,
	opts ...FindOption,
) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.getTimeout())
	defer cancel()

	cursor, err := r.db.
		GetCollection(r.collection).
		Find(
			ctx,
			assertFilter(filter),
			buildOptions(opts),
		)

	if err != nil {
		return nil, logs.Error("error.loading", map[string]any{
			"Subject": r.collection,
			"Error":   err,
		})
	}

	defer cursor.Close(ctx)

	return decodeCursor[T](cursor, r.collection, ctx)
}

func (r *MongoRepository) toCollection[T Document](data T) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.getTimeout())
	defer cancel()

	data.Touch()

	_, err := r.db.
		GetCollection(r.collection).
		ReplaceOne(
			ctx,
			data.Filter(),
			data,
			options.
				Replace().
				SetUpsert(true),
		)

	if err != nil {
		return logs.Error("error.unexpected:save", map[string]any{
			"Error": err,
		})
	}

	return nil
}

// Stream
type Watchable[T any] interface {
	Watch(ctx context.Context) (<-chan T, error)
}

func (r *MongoRepository) watchCollection[T any](ctx context.Context) (<-chan T, error) {
	pipeline := mongodb.Pipeline{
		{{Key: "$match", Value: bson.M{
			"operationType": bson.M{"$in": bson.A{"insert", "replace"}},
		}}},
	}

	stream, err := r.db.GetCollection(r.collection).Watch(ctx, pipeline)
	if err != nil {
		return nil, logs.Error("error.loading", map[string]any{
			"Subject": r.collection,
			"Error":   err,
		})
	}

	out := make(chan T)
	go func() {
		defer close(out)
		defer stream.Close(context.Background())

		for stream.Next(ctx) {
			var event struct {
				FullDocument T `bson:"fullDocument"`
			}
			if err := stream.Decode(&event); err != nil {
				continue
			}

			select {
			case out <- event.FullDocument:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

// Methods options
type FindOption func(*options.FindOptions)

func WithLimit(n int64) FindOption {
	return func(o *options.FindOptions) { o.SetLimit(n) }
}

func WithSort(sort bson.D) FindOption {
	return func(o *options.FindOptions) { o.SetSort(sort) }
}

func WithSkip(n int64) FindOption {
	return func(o *options.FindOptions) { o.SetSkip(n) }
}

// Helpers
func assertFilter(filter bson.M) bson.M {
	if filter == nil {
		return bson.M{}
	}
	return filter
}

func buildOptions(opts []FindOption) *options.FindOptions {
	findOpts := options.Find()
	for _, opt := range opts {
		opt(findOpts)
	}
	return findOpts
}

func decodeCursor[T any](cursor *mongodb.Cursor, collection mango.Collection, ctx context.Context) ([]T, error) {
	var data []T
	if err := cursor.All(ctx, &data); err != nil {
		return nil, logs.Error("error.loading", map[string]any{
			"Subject": collection,
			"Error":   err,
		})
	}

	return data, nil
}
