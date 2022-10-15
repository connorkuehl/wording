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
