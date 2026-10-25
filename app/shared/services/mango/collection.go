package mango

// Collection is a MongoDB collection name.
type Collection string

// String returns the collection name.
func (c Collection) String() string {
	return string(c)
}
