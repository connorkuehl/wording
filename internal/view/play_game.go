package view

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/connorkuehl/wording/internal/wording"
)

//go:embed play_game.tmpl.html
var playGameHTML string

var playGameTmpl = template.Must(template.New("play-game").Parse(playGameHTML))

// PlayGame is the play game page.
type PlayGame struct {
	Token     string
	Length    int
	GameState *wording.GameState
}

// RenderTo renders the play game page.
func (v PlayGame) RenderTo(w io.Writer) error {
	return playGameTmpl.Execute(w, v)
}
