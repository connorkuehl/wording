package wording

import (
	"strings"
	"time"
)

type Game struct {
	AdminToken string
	Token      string
	Answer     string
	ExpiresAt  time.Time
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
