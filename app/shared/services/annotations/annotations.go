package annotations

// NODB opts a command out of the DB, logger and i18n bootstrap.
const (
	NODB = "noDb"
)

// Set returns new annotations with key set to value.
func Set(key, value string) Annotations {
	return Annotations{key: value}
}

// Annotations are cobra command annotations.
type Annotations map[string]string

// Has reports whether key is set, even to an empty value.
func (a Annotations) Has(key string) bool {
	_, ok := a[key]
	return ok
}

// Get returns the value of key, or "".
func (a Annotations) Get(key string) string {
	return a[key]
}

// Set stores value under key and returns a.
func (a Annotations) Set(key, value string) Annotations {
	a[key] = value
	return a
}
