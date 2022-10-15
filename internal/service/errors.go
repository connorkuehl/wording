package service

import "errors"

var (
	ErrGuessLimitReached = errors.New("guess limit reached")
	ErrNotFound          = errors.New("not found")
	ErrCannotContinue    = errors.New("game is over")
)
