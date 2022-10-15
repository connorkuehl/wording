package view

import (
	_ "embed"
	"html/template"
	"io"
)

//go:embed play_game.tmpl.html
var playGameHTML string

var playGameTmpl = template.Must(template.New("play-game").Parse(playGameHTML))

type PlayGame struct {
	Token    string
	Length   int
	Attempts []string
}

func (v PlayGame) RenderTo(w io.Writer) error {
	return playGameTmpl.Execute(w, v)
}
