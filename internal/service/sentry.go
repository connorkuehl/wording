package service

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"

	"github.com/connorkuehl/wording/internal/wording"
)

var uninterestingErrs = []error{
	nil,
	ErrNotFound,
	ErrCannotContinue,
	ErrGuessLimitReached,
}

func isExceptional(err error) bool {
	for _, e := range uninterestingErrs {
		if errors.Is(err, e) {
			return false
		}
	}

	var badInput wording.InputViolations
	return !errors.As(err, &badInput)
}

// Sentry is a wrapper around the service struct that reports exceptional
// errors to sentry.io.
type Sentry struct {
	Service
}

// NewSentry constructs a new Sentry service wrapper.
func NewSentry(svc Service) *Sentry {
	return &Sentry{svc}
}

// CreateGame proxies to Service.CreateGame.
func (s *Sentry) CreateGame(
	ctx context.Context,
	answer string,
	guessLimit int,
) (*wording.Game, error) {
	g, err := s.Service.CreateGame(ctx, answer, guessLimit)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return g, err
}

// Gmae proxies to Service.Game.
func (s *Sentry) Game(ctx context.Context, adminToken string) (*wording.Game, error) {
	g, err := s.Service.Game(ctx, adminToken)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return g, err
}

// GameByToken proxies to Service.GameByToken.
func (s *Sentry) GameByToken(ctx context.Context, token string) (*wording.Game, error) {
	g, err := s.Service.GameByToken(ctx, token)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return g, err
}

// SubmitGuess proxies to Service.SubmitGuess.
func (s *Sentry) SubmitGuess(ctx context.Context, gameToken, playerToken, guess string) error {
	err := s.Service.SubmitGuess(ctx, gameToken, playerToken, guess)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return err
}

// GameState proxies to Service.GameState.
func (s *Sentry) GameState(ctx context.Context, gameToken, playerToken string) (*wording.GameState, error) {
	state, err := s.Service.GameState(ctx, gameToken, playerToken)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return state, err
}

// Plays proxies to Service.Plays.
func (s *Sentry) Plays(ctx context.Context, gameToken, playerToken string) (*wording.Plays, error) {
	plays, err := s.Service.Plays(ctx, gameToken, playerToken)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return plays, err
}

// Stats proxies to Service.Stats.
func (s *Sentry) Stats(ctx context.Context) (wording.Stats, error) {
	stats, err := s.Service.Stats(ctx)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return stats, err
}

// GameStats proxies to Service.GameStats.
func (s *Sentry) GameStats(ctx context.Context, adminToken string) (wording.Stats, error) {
	stats, err := s.Service.GameStats(ctx, adminToken)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return stats, err
}

// DeleteGame proxies to Service.DeleteGame.
func (s *Sentry) DeleteGame(ctx context.Context, adminToken string) error {
	err := s.Service.DeleteGame(ctx, adminToken)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return err
}
