package annotations

// Annotation opting a command into DB + logger + i18n bootstrap.
const (
	NODB = "noDb"
)

func Set(key, value string) Annotations {
	return Annotations{key: value}
}

type Annotations map[string]string

func (a Annotations) Has(key string) bool {
	_, ok := a[key]
	return ok
}

func (a Annotations) Get(key string) string {
	return a[key]
}

func (a Annotations) Set(key, value string) Annotations {
	a[key] = value
	return a
}
