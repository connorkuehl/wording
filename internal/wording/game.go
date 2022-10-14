package wording

import (
	"time"
)

type Game struct {
	AdminToken string
	Token      string
	Answer     string
	ExpiresAt  time.Time
	GuessLimit int
}
