package wording

import (
	"errors"
	"fmt"
	"strings"
)

type Game struct {
	AdminToken string
	Token      string
	Answer     string
	GuessLimit int
}

type Character struct {
	Value     string
	IsCorrect bool
	IsPartial bool
}

type Attempt []Character

func (a Attempt) String() string {
	var s strings.Builder

	for _, ch := range a {
		s.WriteString(ch.Value)
	}

	return s.String()
}

type GameState struct {
	Attempts     []Attempt
	CanContinue  bool
	IsVictorious bool
	GameOver     bool
}

func Evaluate(answer, guess string) Attempt {
	// TODO: consider making a type that maintains the
	// invariant that both answer and guess must be same
	// length

	key := make(map[rune][]int)
	for i, r := range answer {
		key[r] = append(key[r], i)
	}

	var at Attempt
	for i, r := range guess {
		ch := Character{
			Value: string(r),
		}

		locs, ok := key[r]
		for _, l := range locs {
			if l == i {
				ch.IsCorrect = true
				break
			}
		}
		if ok && !ch.IsCorrect {
			ch.IsPartial = true
		}

		at = append(at, ch)
	}

	return at
}

type InputViolations map[string][]error

func (i InputViolations) Error() string {
	return fmt.Sprintf("%v", map[string][]error(i))
}

func ValidateGuessLimit(limit int) error {
	violations := make(InputViolations)

	if limit < 1 || limit > 16 {
		violations["guesses allowed"] = append(violations["guesses allowed"], errors.New("must be between 1-16 characters long"))
	}

	if len(violations) > 0 {
		return violations
	}

	return nil
}

func ValidateAnswer(answer string) error {
	violations := make(InputViolations)

	if !isAlpha(answer) {
		violations["answer"] = append(violations["answer"], errors.New("has non-alphabetical characters"))
	}

	if len(violations) > 0 {
		return violations
	}

	return nil
}

func ValidateGuess(guess, answer string, previousGuesses []string) error {
	violations := make(InputViolations)

	if len(guess) != len(answer) {
		violations["guess"] = append(violations["guess"], fmt.Errorf("guess must be %d characters long", len(answer)))
	}

	if !isAlpha(guess) {
		violations["guess"] = append(violations["guess"], errors.New("has non-alphabetical characters"))
	}

	for _, g := range previousGuesses {
		if g == guess {
			violations["guess"] = append(violations["guess"], errors.New("has already been tried"))
		}
	}

	if len(violations) > 0 {
		return violations
	}

	return nil
}

func isAlpha(s string) bool {
	legalValues := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, r := range s {
		if !strings.ContainsRune(legalValues, r) {
			return false
		}
	}
	return true
}
