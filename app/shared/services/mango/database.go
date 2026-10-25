package mango

var db *Mongo

// GetDb returns the shared connection, nil before SetDb.
func GetDb() *Mongo {
	return db
}

// SetDb sets the shared connection.
func SetDb(connection *Mongo) {
	db = connection
}
