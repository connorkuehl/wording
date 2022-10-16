package view

import (
	_ "embed"
	"html/template"
	"io"
)

//go:embed manage_game.tmpl.html
var manageGameHTML string

var manageGameTmpl = template.Must(template.New("manage-game").Parse(manageGameHTML))

type ManageGame struct {
	BaseURL        string
	AdminToken     string
	Token          string
	Answer         string
	GuessesAllowed int
}

func (m ManageGame) RenderTo(w io.Writer) error {
	return manageGameTmpl.Execute(w, m)
}
