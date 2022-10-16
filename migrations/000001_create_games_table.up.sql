CREATE TABLE IF NOT EXISTS games (
    admin_token TEXT PRIMARY KEY,
    token TEXT,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    answer TEXT NOT NULL,
    guess_limit INTEGER NOT NULL
);