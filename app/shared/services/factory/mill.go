package factory

import (
	plumprimitive "github.com/cyril-fl/plum/primitive"
)

type Mill struct {
	plumprimitive.Fallible
}

/*
	NOTE Snippet - Mill

	type Mill = mill

	type mill struct {
		factory.Mill

	}

	type Template struct {

	}

	func New(params Template) *mill {
		return &mill{

		}
	}

	// Getters & Setters

	// Handlers

	// Helpers

*/
