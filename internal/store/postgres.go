package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"

	"github.com/connorkuehl/wording/internal/wording"
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

func (s *PostgresStore) CreateGame(ctx context.Context, adminToken, token, answer string, guessLimit int, expiresAt time.Time) (*wording.Game, error) {
	query := `
	INSERT INTO games (
		admin_token,
		token,
		expires_at,
		answer,
		guess_limit
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)
	`

	_, err := s.db.ExecContext(ctx, query, adminToken, token, expiresAt, answer, guessLimit)
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

func (s *PostgresStore) Game(ctx context.Context, adminToken string) (*wording.Game, error) {
	query := `SELECT token, expires_at, answer, guess_limit FROM games WHERE admin_token = $1`

	var game wording.Game
	err := s.db.QueryRowContext(ctx, query, adminToken).
		Scan(&game.Token, &game.ExpiresAt, &game.Answer, &game.GuessLimit)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	game.AdminToken = adminToken

	return &game, nil
}
