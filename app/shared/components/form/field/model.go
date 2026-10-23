package field

import (
	"errors"
	"solopg/app/shared/services/i19n"

	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// - Field - //
type field[T any] struct {
	id           string
	kind         kind
	label        string
	focus        bool
	defaultvalue T
	validator    func(T) error
	required     bool
	nolabel      bool
	err          []error
}

type kind string

type Template[T any] struct {
	ID           string
	Label        string
	Defaultvalue T
	Validator    func(T) error
	Required     bool
	NoLabel      bool
}

const (
	// FieldTypeInput represents a text input field.
	INPUT kind = "input"
	// FieldTypeSelect represents a select field.
	SELECT kind = "select"
)

// newField creates a form field with an automatically generated identifier.
func newField[T any](id, label string, kind kind, value T, validator func(T) error, required, nolabel bool) field[T] {
	return field[T]{
		id:           id,
		kind:         kind,
		label:        label,
		defaultvalue: value,
		validator:    validator,
		required:     required,
		nolabel:      nolabel,
	}
}

// - Methods --//
func (f *field[T]) ID() string {
	return f.id
}

func (f *field[T]) Label() string {
	label := i19n.Localize(f.label)
	if label == "" {
		label = f.label
	}
	if label == "" {
		label = f.id
	}
	return plumtransform.Capitalize(label)
}

func (f *field[T]) SetError(err error) {
	f.err = append(f.err, err)
}

func (f *field[T]) GetError() error {
	if len(f.err) > 0 {
		return errors.Join(f.err...)
	}

	return nil
}

func (f *field[T]) ResetError() {
	f.err = nil
}

// - Helper --//
func (f *field[T]) validate(value T) error {
	if f.validator != nil {
		err := f.validator(value)

		f.SetError(err)
		return err
	}

	return nil
}

func (f *field[T]) IsFocused() bool {
	return f.focus
}
