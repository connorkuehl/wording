package service

import "errors"

var (
	ErrGuessLimitReached = errors.New("guess limit reached")
	ErrNotFound          = errors.New("not found")
)
