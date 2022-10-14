package service

import (
	"context"
	"time"

	"github.com/connorkuehl/wording/internal/wording"
)

var now = time.Now

//go:generate mockery --name Store --case underscore --with-expecter --testonly --inpackage
type Store interface {
	CreateGame(ctx context.Context, adminToken, token, answer string, guessLimit int, expiresAt time.Time) (*wording.Game, error)
	Game(ctx context.Context, adminToken string) (*wording.Game, error)
	GameByToken(ctx context.Context, token string) (*wording.Game, error)
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
