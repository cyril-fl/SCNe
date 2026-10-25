package field

import (
	"errors"
	"scne/app/shared/services/i19n"

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

// Template holds the options shared by every field kind.
type Template[T any] struct {
	ID           string
	Label        string
	Defaultvalue T
	Validator    func(T) error
	Required     bool
	NoLabel      bool
}

const (
	// INPUT is a text input field.
	INPUT kind = "input"
	// SELECT is a select field.
	SELECT kind = "select"
	// TEXTAREA is a multi-line text field.
	TEXTAREA kind = "textarea"
)

// newField creates the base state shared by every field kind.
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
