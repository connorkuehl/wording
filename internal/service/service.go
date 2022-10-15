package service

import (
	"context"
	"errors"
	"time"

	"github.com/connorkuehl/wording/internal/store"
	"github.com/connorkuehl/wording/internal/wording"
)

var now = time.Now

//go:generate mockery --name Store --case underscore --with-expecter --testonly --inpackage
type Store interface {
	CreateGame(ctx context.Context, adminToken, token, answer string, guessLimit int, expiresAt time.Time) (*wording.Game, error)
	Game(ctx context.Context, adminToken string) (*wording.Game, error)
	GameByToken(ctx context.Context, token string) (*wording.Game, error)
	Plays(ctx context.Context, gameToken, playerToken string) (*wording.Plays, error)
	PutPlays(ctx context.Context, gameToken, playerToken string, plays *wording.Plays) error
}

//go:generate mockery --name TokenGenerator --case underscore --with-expecter --testonly --inpackage
type TokenGenerator interface {
	NewToken() string
}

type Service struct {
	store               Store
	adminTokenGenerator TokenGenerator
	gameTokenGenerator  TokenGenerator
}

func New(store Store, adminTokenGenerator, gameTokenGenerator TokenGenerator) *Service {
	return &Service{
		store:               store,
		adminTokenGenerator: adminTokenGenerator,
		gameTokenGenerator:  gameTokenGenerator,
	}
}

func (s *Service) CreateGame(
	ctx context.Context,
	answer string,
	guessLimit int,
	expiresAfter time.Duration,
) (*wording.Game, error) {
	expiresAt := now().Add(expiresAfter)
	return s.store.CreateGame(ctx, s.adminTokenGenerator.NewToken(), s.gameTokenGenerator.NewToken(), answer, guessLimit, expiresAt)
}

func (s *Service) Game(ctx context.Context, adminToken string) (*wording.Game, error) {
	return s.store.Game(ctx, adminToken)
}

func (s *Service) GameByToken(ctx context.Context, token string) (*wording.Game, error) {
	return s.store.GameByToken(ctx, token)
}

func (s *Service) SubmitGuess(ctx context.Context, gameToken, playerToken, guess string) error {
	game, err := s.store.GameByToken(ctx, gameToken)
	if err != nil {
		return err
	}

	plays, err := s.store.Plays(ctx, gameToken, playerToken)
	if errors.Is(err, store.ErrNotFound) {
		plays = &wording.Plays{}
		err = nil
	}
	if err != nil {
		return err
	}

	if len(plays.Attempts) >= game.GuessLimit {
		return ErrGuessLimitReached
	}

	state := plays.Evaluate(game.Answer, game.GuessLimit)
	if !state.CanContinue {
		return ErrCannotContinue
	}

	plays.Attempts = append(plays.Attempts, guess)
	err = s.store.PutPlays(ctx, gameToken, playerToken, plays)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) NewPlayerToken(ctx context.Context) string {
	return s.adminTokenGenerator.NewToken()
}

func (s *Service) GameState(ctx context.Context, gameToken, playerToken string) (*wording.GameState, error) {
	game, err := s.store.GameByToken(ctx, gameToken)
	if errors.Is(err, store.ErrNotFound) {
		return nil, ErrNotFound
	}

	plays, err := s.store.Plays(ctx, gameToken, playerToken)
	if errors.Is(err, store.ErrNotFound) {
		return &wording.GameState{CanContinue: true}, nil
	}

	return plays.Evaluate(game.Answer, game.GuessLimit), nil
}

func (s *Service) Plays(ctx context.Context, gameToken, playerToken string) (*wording.Plays, error) {
	plays, err := s.store.Plays(ctx, gameToken, playerToken)
	if errors.Is(err, store.ErrNotFound) {
		return new(wording.Plays), nil
	}
	if err != nil {
		return nil, err
	}

	return plays, nil
}
