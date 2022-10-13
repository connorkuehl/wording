package service

import (
	"context"
	"strings"
	"time"

	"github.com/connorkuehl/wording/internal/wording"
)

type Nower func() time.Time

//go:generate mockery --name Store --case underscore --with-expecter --testonly --inpackage
type Store interface {
	CreateGame(ctx context.Context, adminToken, answer string, guessLimit int, expiresAt time.Time) (*wording.Game, error)
}

//go:generate mockery --name TokenGenerator --case underscore --with-expecter --testonly --inpackage
type TokenGenerator interface {
	Adjective() string
	Noun() string
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
	adminToken := strings.Join([]string{s.tokGen.Adjective(), s.tokGen.Noun()}, "-")

	return s.store.CreateGame(ctx, adminToken, answer, guessLimit, expiresAt)
}
