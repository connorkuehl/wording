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
	if errors.As(err, &badInput) {
		return false
	}

	return true
}

type Sentry struct {
	Service
}

func NewSentry(svc Service) *Sentry {
	return &Sentry{svc}
}

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

func (s *Sentry) Game(ctx context.Context, adminToken string) (*wording.Game, error) {
	g, err := s.Service.Game(ctx, adminToken)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return g, err
}

func (s *Sentry) GameByToken(ctx context.Context, token string) (*wording.Game, error) {
	g, err := s.Service.GameByToken(ctx, token)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return g, err
}

func (s *Sentry) SubmitGuess(ctx context.Context, gameToken, playerToken, guess string) error {
	err := s.Service.SubmitGuess(ctx, gameToken, playerToken, guess)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return err
}

func (s *Sentry) GameState(ctx context.Context, gameToken, playerToken string) (*wording.GameState, error) {
	state, err := s.Service.GameState(ctx, gameToken, playerToken)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return state, err
}

func (s *Sentry) Plays(ctx context.Context, gameToken, playerToken string) (*wording.Plays, error) {
	plays, err := s.Service.Plays(ctx, gameToken, playerToken)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return plays, err
}

func (s *Sentry) Stats(ctx context.Context) (wording.Stats, error) {
	stats, err := s.Service.Stats(ctx)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return stats, err
}

func (s *Sentry) GameStats(ctx context.Context, adminToken string) (wording.Stats, error) {
	stats, err := s.Service.GameStats(ctx, adminToken)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return stats, err
}

func (s *Sentry) DeleteGame(ctx context.Context, adminToken string) error {
	err := s.Service.DeleteGame(ctx, adminToken)
	if isExceptional(err) {
		sentry.CaptureException(err)
	}
	return err
}
