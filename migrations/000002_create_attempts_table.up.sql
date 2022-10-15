CREATE TABLE IF NOT EXISTS attempts (
    game_token TEXT NOT NULL,
    player_token TEXT NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    guesses text[] NOT NULL,
    UNIQUE (game_token, player_token)
);