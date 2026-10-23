package mango

import (
	"context"

	"time"

	plumprimitive "github.com/cyril-fl/plum/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var env = NewEnvironment()

// -- Mongo & Database -- //
type Mongo struct {
	plumprimitive.Fallible

	client   *mongo.Client
	instance *mongo.Database
}

func New() *Mongo {
	return &Mongo{}
}

// Getters & Setters
func (db *Mongo) GetCollection(collection Collection) *mongo.Collection {
	return db.instance.Collection(collection.String())
}

func (db *Mongo) setClient(client *mongo.Client) {
	db.client = client
}

func (db *Mongo) setInstance(client *mongo.Client) {
	db.instance = client.Database(env.GetEnvDBName())
}

// Methods
func (db *Mongo) Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.
		Connect(ctx, options.Client().
			ApplyURI(env.GetEnvURI()).
			SetDirect(true))

	if err != nil {
		db.SetErr(err)
		return
	}

	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx)
		db.SetErr(err)
		return
	}

	db.setClient(client)
	db.setInstance(client)
}

func (db *Mongo) Disconnect() {
	if db.client == nil {
		return
	}
	db.client.Disconnect(context.Background())
}

// Helpers
func Init() (db *Mongo, err error) {
	env.ReadEnv()

	db = New()
	db.Connect()
	if db.HasErr() {
		return nil, db.GetErr()
	}
	return db, nil
}
