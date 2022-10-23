package wording

// LifetimeScope is the name of the row where global game stats are stored.
const LifetimeScope = "lifetime"

// Stats are some of the interesting stats to show players.
type Stats struct {
	GamesCreated int
	GamesWon     int
	GuessesMade  int
}

// IncrementStats is just a type-name wrapper suggesting that each stat's field
// is being incremented and is not an absolute value.
type IncrementStats struct {
	Stats
}
