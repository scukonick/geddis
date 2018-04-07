package db

import "errors"

var (
	// ErrNotFound is returned when the element
	// with needed key is not presented in the store
	ErrNotFound = errors.New("element not found")

	// ErrInvalidType is returned when client tries to
	// get type which is not store by  the key
	// or does operations on type which does not support them
	ErrInvalidType = errors.New("invalid store type")
)
