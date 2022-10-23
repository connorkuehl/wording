package wording

// Plays are all of the guesses a player has made.
type Plays struct {
	Attempts []string
}

// Evaluate checks all of the player's attempts and produces a GameState
// snapshot.
func (p *Plays) Evaluate(answer string, guessLimit int) *GameState {
	var ats []Attempt
	for _, play := range p.Attempts {
		ats = append(ats, Evaluate(answer, play))
	}

	state := GameState{
		Attempts: ats,
	}

	for _, attempt := range state.Attempts {
		correct := true

		for _, ch := range attempt {
			correct = correct && ch.IsCorrect
		}

		if correct {
			state.IsVictorious = true
			break
		}
	}

	state.GameOver = len(state.Attempts) >= guessLimit
	state.CanContinue = !state.IsVictorious && !state.GameOver

	return &state
}
