package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/connorkuehl/wording/internal/store"
	"github.com/connorkuehl/wording/internal/wording"
)

//go:generate mockery --name Store --case underscore --with-expecter --testonly --inpackage
type Store interface {
	CreateGame(ctx context.Context, adminToken, token, answer string, guessLimit int) (*wording.Game, error)
	Game(ctx context.Context, adminToken string) (*wording.Game, error)
	GameByToken(ctx context.Context, token string) (*wording.Game, error)
	Plays(ctx context.Context, gameToken, playerToken string) (*wording.Plays, error)
	PutPlays(ctx context.Context, gameToken, playerToken string, plays *wording.Plays) error
	IncrementStats(ctx context.Context, stats wording.IncrementStats) error
	GameStats(ctx context.Context, adminToken string) (wording.Stats, error)
	Stats(ctx context.Context) (wording.Stats, error)
	DeleteGame(ctx context.Context, adminToken string) error
}

//go:generate mockery --name TokenGenerator --case underscore --with-expecter --testonly --inpackage
type TokenGenerator interface {
	NewToken() string
}

type Service interface {
	CreateGame(ctx context.Context, answer string, guessLimit int) (*wording.Game, error)
	DeleteGame(ctx context.Context, adminToken string) error
	Game(ctx context.Context, adminToken string) (*wording.Game, error)
	GameByToken(ctx context.Context, token string) (*wording.Game, error)
	GameState(ctx context.Context, gameToken, playerToken string) (*wording.GameState, error)
	GameStats(ctx context.Context, adminToken string) (wording.Stats, error)
	NewPlayerToken(ctx context.Context) string
	Plays(ctx context.Context, gameToken, playerToken string) (*wording.Plays, error)
	Stats(ctx context.Context) (wording.Stats, error)
	SubmitGuess(ctx context.Context, gameToken, playerToken, guess string) error
}

type service struct {
	store               Store
	adminTokenGenerator TokenGenerator
	gameTokenGenerator  TokenGenerator
}

// New creates a new service.
func New(store Store, adminTokenGenerator, gameTokenGenerator TokenGenerator) *service {
	return &service{
		store:               store,
		adminTokenGenerator: adminTokenGenerator,
		gameTokenGenerator:  gameTokenGenerator,
	}
}

// CreateGame creates a new guess-the-word game.
func (s *service) CreateGame(
	ctx context.Context,
	answer string,
	guessLimit int,
) (*wording.Game, error) {
	err := wording.ValidateAnswer(answer)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	answer = strings.ToLower(answer)

	err = wording.ValidateGuessLimit(guessLimit)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	game, err := s.store.CreateGame(ctx, s.adminTokenGenerator.NewToken(), s.gameTokenGenerator.NewToken(), answer, guessLimit)
	if err != nil {
		return nil, err
	}

	err = s.store.IncrementStats(ctx, wording.IncrementStats{Stats: wording.Stats{GamesCreated: 1}})
	if err != nil {
		// TODO
		log.Println("increment:", err)
	}

	return game, nil
}

// Game fetches the game identified by the adminToken.
func (s *service) Game(ctx context.Context, adminToken string) (*wording.Game, error) {
	game, err := s.store.Game(ctx, adminToken)
	if errors.Is(err, store.ErrNotFound) {
		err = ErrNotFound
	}
	return game, err
}

// GameByToken fetches the game identified by token.
func (s *service) GameByToken(ctx context.Context, token string) (*wording.Game, error) {
	game, err := s.store.GameByToken(ctx, token)
	if errors.Is(err, store.ErrNotFound) {
		err = ErrNotFound
	}
	return game, err
}

// SubmitGuess records the player's guess.
func (s *service) SubmitGuess(ctx context.Context, gameToken, playerToken, guess string) error {
	guess = strings.ToLower(guess)

	game, err := s.store.GameByToken(ctx, gameToken)
	if errors.Is(err, store.ErrNotFound) {
		err = ErrNotFound
	}
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

	err = wording.ValidateGuess(guess, game.Answer, plays.Attempts)
	if err != nil {
		return fmt.Errorf("invalid input: %w", err)
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

	plays, err = s.store.Plays(ctx, gameToken, playerToken)
	if err != nil {
		return err
	}
	state = plays.Evaluate(game.Answer, game.GuessLimit)

	incWins := 0
	if state.IsVictorious {
		incWins = 1
	}

	err = s.store.IncrementStats(ctx, wording.IncrementStats{Stats: wording.Stats{GuessesMade: 1, GamesWon: incWins}})
	if err != nil {
		// TODO
		log.Println("increment:", err)
	}

	return nil
}

// NewPlayerToken allocates a new player token.
func (s *service) NewPlayerToken(ctx context.Context) string {
	return s.adminTokenGenerator.NewToken()
}

// GameState returns a snapshot of a player's progress against a given game.
func (s *service) GameState(ctx context.Context, gameToken, playerToken string) (*wording.GameState, error) {
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

// Plays fetches a player's attempts against a game.
func (s *service) Plays(ctx context.Context, gameToken, playerToken string) (*wording.Plays, error) {
	plays, err := s.store.Plays(ctx, gameToken, playerToken)
	if errors.Is(err, store.ErrNotFound) {
		return new(wording.Plays), nil
	}
	if err != nil {
		return nil, err
	}

	return plays, nil
}

// Stats fetches the application's lifetime stats.
func (s *service) Stats(ctx context.Context) (wording.Stats, error) {
	return s.store.Stats(ctx)
}

// DeleteGame deletes the specified game and all of the attempts made against it.
func (s *service) DeleteGame(ctx context.Context, adminToken string) error {
	err := s.store.DeleteGame(ctx, adminToken)
	if errors.Is(err, store.ErrNotFound) {
		err = ErrNotFound
	}
	return err
}

// GameStats returns a specific game's stats.
func (s *service) GameStats(ctx context.Context, adminToken string) (wording.Stats, error) {
	return s.store.GameStats(ctx, adminToken)
}
