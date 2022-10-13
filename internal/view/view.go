package view

import (
	_ "embed"
)

//go:embed home.html
var homePage string

func Home() string {
	return homePage
}
