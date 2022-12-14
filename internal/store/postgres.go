package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"github.com/connorkuehl/wording/internal/wording"
)

// PostgresStore is a PostgreSQL-backed persistence layer for the game.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore opens a connection to the Postgres store.
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

// Close closes the connection to the Postgres store.
func (s *PostgresStore) Close() error {
	return s.db.Close()
}

// CreateGame creates a game.
func (s *PostgresStore) CreateGame(ctx context.Context, adminToken, token, answer string, guessLimit int) (*wording.Game, error) {
	query := `
	INSERT INTO games (
		admin_token,
		token,
		answer,
		guess_limit
	) VALUES (
		$1,
		$2,
		$3,
		$4
	)
	`

	_, err := s.db.ExecContext(ctx, query, adminToken, token, answer, guessLimit)
	if err != nil {
		return nil, err
	}

	game := &wording.Game{
		AdminToken: adminToken,
		Token:      token,
		Answer:     answer,
		GuessLimit: guessLimit,
	}

	return game, nil
}

// Game fetches a game.
func (s *PostgresStore) Game(ctx context.Context, adminToken string) (*wording.Game, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	query := `SELECT token, answer, guess_limit FROM games WHERE admin_token = $1`

	game := wording.Game{AdminToken: adminToken}
	err = tx.QueryRowContext(ctx, query, adminToken).
		Scan(&game.Token, &game.Answer, &game.GuessLimit)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `UPDATE games SET accessed_at = NOW() WHERE admin_token = $1`, adminToken)
	if err != nil {
		return nil, err
	}

	return &game, tx.Commit()
}

// GameByToken fetches a game by the the specified token.
func (s *PostgresStore) GameByToken(ctx context.Context, token string) (*wording.Game, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	query := `SELECT admin_token, answer, guess_limit FROM games WHERE token = $1`

	game := wording.Game{
		Token: token,
	}
	err = tx.QueryRowContext(ctx, query, token).
		Scan(&game.AdminToken, &game.Answer, &game.GuessLimit)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `UPDATE games SET accessed_at = NOW() WHERE token = $1`, token)
	if err != nil {
		return nil, err
	}

	return &game, tx.Commit()
}

// Plays fetches a player's attempts against a given game.
func (s *PostgresStore) Plays(ctx context.Context, gameToken, playerToken string) (*wording.Plays, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	plays := &wording.Plays{}

	query := `SELECT guesses FROM attempts WHERE game_token=$1 AND player_token=$2`
	err = tx.QueryRowContext(ctx, query, gameToken, playerToken).Scan(pq.Array(&plays.Attempts))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	_, err = tx.ExecContext(ctx, `UPDATE games SET accessed_at = NOW() WHERE token = $1`, gameToken)
	if err != nil {
		return nil, err
	}

	return plays, tx.Commit()
}

// PutPlays updates a player's attempts against a game.
func (s *PostgresStore) PutPlays(ctx context.Context, gameToken, playerToken string, plays *wording.Plays) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

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

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	query = `UPDATE games SET accessed_at = NOW(),
							  modified_at = NOW()
							  WHERE token = $1`

	_, err = tx.ExecContext(ctx, query, gameToken)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// IncrementStats adjusts overall stats for the application.
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

// Stats fetches the overall lifetime stats.
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

// GameStats fetches stats for an individual game.
func (s *PostgresStore) GameStats(ctx context.Context, adminToken string) (wording.Stats, error) {
	var stats wording.Stats

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return stats, err
	}
	defer func() { _ = tx.Rollback() }()

	query := `SELECT COUNT(*)
	FROM attempts WHERE
	game_token = (SELECT token FROM games WHERE admin_token = $1) AND
	(SELECT answer FROM games WHERE admin_token = $1) = ANY (guesses);`
	err = tx.QueryRowContext(ctx, query, adminToken).Scan(&stats.GamesWon)
	if err != nil {
		return stats, err
	}
	query = `SELECT array_length(guesses, 1) FROM attempts WHERE game_token = (SELECT token FROM games WHERE admin_token = $1)`
	rows, err := tx.QueryContext(ctx, query, adminToken)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var playerGuesses int
		err := rows.Scan(&playerGuesses)
		if err != nil {
			return stats, err
		}

		stats.GuessesMade += playerGuesses
	}

	return stats, tx.Commit()
}

// DeleteGame deletes the game and all of the attempts recorded against it.
func (s *PostgresStore) DeleteGame(ctx context.Context, adminToken string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `DELETE FROM attempts WHERE game_token = (SELECT token FROM games WHERE admin_token = $1)`, adminToken)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM games WHERE admin_token = $1`, adminToken)
	if err != nil {
		return err
	}

	return tx.Commit()
}
