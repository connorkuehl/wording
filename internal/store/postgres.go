package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/connorkuehl/wording/internal/wording"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: retry
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	s := &PostgresStore{
		db: db,
	}

	return s, nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

func (s *PostgresStore) CreateGame(ctx context.Context, adminToken, answer string, guessLimit int, expiresAt time.Time) (*wording.Game, error) {
	query := `
	INSERT INTO games (
		admin_token,
		expires_at,
		answer,
		guess_limit
	) VALUES (
		$1,
		$2,
		$3,
		$4
	)
	`

	_, err := s.db.ExecContext(ctx, query, adminToken, expiresAt, answer, guessLimit)
	if err != nil {
		return nil, err
	}

	game := &wording.Game{
		AdminToken: adminToken,
		ExpiresAt:  expiresAt,
		Answer:     answer,
		GuessLimit: guessLimit,
	}

	return game, nil
}
