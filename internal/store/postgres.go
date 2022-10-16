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

func (s *PostgresStore) Game(ctx context.Context, adminToken string) (*wording.Game, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

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

func (s *PostgresStore) GameByToken(ctx context.Context, token string) (*wording.Game, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

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

func (s *PostgresStore) Plays(ctx context.Context, gameToken, playerToken string) (*wording.Plays, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

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

func (s *PostgresStore) PutPlays(ctx context.Context, gameToken, playerToken string, plays *wording.Plays) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

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

	_, err = tx.ExecContext(ctx, `UPDATE games SET accessed_at = NOW(), modified_at = NOW() WHERE token = $1`, gameToken)
	if err != nil {
		return err
	}

	return tx.Commit()
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

func (s *PostgresStore) DeleteGame(ctx context.Context, adminToken string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

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
