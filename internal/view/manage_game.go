package view

import (
	_ "embed"
	"html/template"
	"io"
)

//go:embed manage_game.tmpl.html
var manageGameHTML string

var manageGameTmpl = template.Must(template.New("manage-game").Parse(manageGameHTML))

// ManageGame is the game management screen.
type ManageGame struct {
	BaseURL        string
	AdminToken     string
	Token          string
	Answer         string
	GuessesAllowed int
	GuessesMade    int
	CorrectGuesses int
}

// RenderTo renders the management page.
func (m ManageGame) RenderTo(w io.Writer) error {
	return manageGameTmpl.Execute(w, m)
}
