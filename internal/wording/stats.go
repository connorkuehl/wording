package wording

const LifetimeScope = "lifetime"

type Stats struct {
	GamesCreated int
	GamesWon     int
	GuessesMade  int
}

type IncrementStats struct {
	Stats
}
