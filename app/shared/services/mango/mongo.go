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

// Mongo is a MongoDB connection to the app database.
type Mongo struct {
	plumprimitive.Fallible

	client   *mongo.Client
	instance *mongo.Database
}

// New returns an unconnected Mongo.
func New() *Mongo {
	return &Mongo{}
}

// Getters & Setters

// GetCollection returns collection from the app database.
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

// Connect connects and pings within 10 seconds, storing any error.
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

// Disconnect closes the connection, if any.
func (db *Mongo) Disconnect() {
	if db.client == nil {
		return
	}
	db.client.Disconnect(context.Background())
}

// Helpers

// Init reads the environment and returns a connected Mongo.
func Init() (db *Mongo, err error) {
	env.ReadEnv()

	db = New()
	db.Connect()
	if db.HasErr() {
		return nil, db.GetErr()
	}
	return db, nil
}
