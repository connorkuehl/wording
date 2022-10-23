package store

import "errors"

var (
	// ErrNotFound indicates the record could not be found.
	ErrNotFound = errors.New("not found")
)
