package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

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

func (s *PostgresStore) GameByToken(ctx context.Context, token string) (*wording.Game, error) {
	query := `SELECT admin_token, expires_at, answer, guess_limit FROM games WHERE token = $1`

	game := wording.Game{
		Token: token,
	}
	err := s.db.QueryRowContext(ctx, query, token).
		Scan(&game.AdminToken, &game.ExpiresAt, &game.Answer, &game.GuessLimit)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (s *PostgresStore) Plays(ctx context.Context, gameToken, playerToken string) (*wording.Plays, error) {
	plays := &wording.Plays{}

	query := `SELECT guesses FROM attempts WHERE game_token=$1 AND player_token=$2`
	err := s.db.QueryRowContext(ctx, query, gameToken, playerToken).Scan(pq.Array(&plays.Attempts))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	return plays, nil
}

func (s *PostgresStore) PutPlays(ctx context.Context, gameToken, playerToken string, plays *wording.Plays) error {
	query := `INSERT INTO attempts (
		game_token,
		player_token,
		guesses
	) VALUES (
		$1,
		$2,
		$3
	) ON CONFLICT (game_token, player_token) DO UPDATE SET guesses = $3`
	args := []any{gameToken, playerToken, pq.Array(plays.Attempts)}

	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *PostgresStore) IncrementStats(ctx context.Context, stats wording.IncrementStats) error {
	query := `INSERT INTO stats (
		scope,
		games_created,
		games_won,
		guesses_made
	) VALUES (
		$1, $2, $3, $4
	) ON CONFLICT (scope) DO UPDATE SET
	games_created = stats.games_created + $2,
	games_won = stats.games_won + $3,
	guesses_made = stats.guesses_made + $4`
	args := []any{wording.LifetimeScope, stats.GamesCreated, stats.GamesWon, stats.GuessesMade}

	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *PostgresStore) Stats(ctx context.Context) (wording.Stats, error) {
	var stats wording.Stats

	query := `SELECT games_created, games_won, guesses_made FROM stats WHERE scope = $1`
	err := s.db.QueryRowContext(ctx, query, wording.LifetimeScope).
		Scan(&stats.GamesCreated, &stats.GamesWon, &stats.GuessesMade)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	if err != nil {
		return wording.Stats{}, err
	}

	return stats, nil
}
