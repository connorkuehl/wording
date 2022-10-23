package view

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/connorkuehl/wording/internal/wording"
)

//go:embed home.tmpl.html
var homePage string

var homeTmpl = template.Must(template.New("home").Parse(homePage))

// Home is the homepage/game creation screen.
type Home struct {
	Stats wording.Stats
}

// RenderTo renders the home page.
func (h Home) RenderTo(w io.Writer) error {
	return homeTmpl.Execute(w, h)
}
