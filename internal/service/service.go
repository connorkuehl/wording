package service

import (
	"context"
	"time"

	"github.com/connorkuehl/wording/internal/wording"
)

type Nower func() time.Time

//go:generate mockery --name Store --case underscore --with-expecter --testonly --inpackage
type Store interface {
	CreateGame(ctx context.Context, adminToken, answer string, guessLimit int, expiresAt time.Time) (*wording.Game, error)
	Game(ctx context.Context, adminToken string) (*wording.Game, error)
}

//go:generate mockery --name TokenGenerator --case underscore --with-expecter --testonly --inpackage
type TokenGenerator interface {
	NewToken() string
}

type Service struct {
	store  Store
	tokGen TokenGenerator
	now    Nower
}

func New(store Store, tokGen TokenGenerator, now Nower) *Service {
	return &Service{
		store:  store,
		tokGen: tokGen,
		now:    now,
	}
}

func (s *Service) CreateGame(
	ctx context.Context,
	answer string,
	guessLimit int,
	expiresAt time.Time,
) (*wording.Game, error) {
	return s.store.CreateGame(ctx, s.tokGen.NewToken(), answer, guessLimit, expiresAt)
}

}
