package wording

import (
	"time"
)

type Game struct {
	AdminToken string
	Answer     string
	ExpiresAt  time.Time
	GuessLimit int
}
