CREATE TABLE IF NOT EXISTS stats (
    scope text PRIMARY KEY,
    games_created int NOT NULL,
    games_won int NOT NULL,
    guesses_made int NOT NULL
);