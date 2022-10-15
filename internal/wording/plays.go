package wording

type Plays struct {
	Attempts []string
}

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
	state.CanContinue = !state.IsVictorious || state.GameOver

	return &state
}
