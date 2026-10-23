package mango

var db *Mongo

func GetDb() *Mongo {
	return db
}

func SetDb(connection *Mongo) {
	db = connection
}
