package db

import "errors"

var (
	// ErrNotFound is returned when the element
	// with needed key is not presented in the store
	ErrNotFound = errors.New("element not found")
)
