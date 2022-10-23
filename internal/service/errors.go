package service

import "errors"

var (
	// ErrGuessLimitReached indicates a player has no more guesses
	// remaining in the game.
	ErrGuessLimitReached = errors.New("guess limit reached")

	// ErrNotFound means the resource does not exist.
	ErrNotFound = errors.New("not found")

	// ErrCannotContinue indicates the player has no more guesses
	// or they have already won.
	ErrCannotContinue = errors.New("game is over")
)
