package generator

import (
	"github.com/connorkuehl/wording/internal/randword"

	log "github.com/sirupsen/logrus"
)

// Tokener produces a token.
type Tokener interface {
	NewToken() string
}

// FallibleTokener tries to produce a token.
type FallibleTokener interface {
	NewToken() (string, error)
}

// FallibleGenerator tries to produce a token with a mechanism that
// might fail, and if so, it will fallback to an infallible mechanism.
type FallibleGenerator struct {
	try      FallibleTokener
	fallback Tokener
}

// NewFallibleGenerator creates a new FallibleGenerator.
func NewFallibleGenerator(try FallibleTokener, fallback Tokener) *FallibleGenerator {
	return &FallibleGenerator{
		try:      try,
		fallback: fallback,
	}
}

// NewToken creates a token. If the fallible tokener fails, the
// fallback is used instead.
func (g *FallibleGenerator) NewToken() string {
	tok, err := g.try.NewToken()
	if err != nil {
		log.WithError(err).Warn("using fallback token generator")
		tok = g.fallback.NewToken()
	}
	return tok
}

// HumanReadable returns a URL slug full of human readable words.
type HumanReadable struct {
	client *randword.Client
}

// NewHumanReadable returns a human-readable token/slug generator.
func NewHumanReadable(client *randword.Client) *HumanReadable {
	return &HumanReadable{
		client: client,
	}
}

// NewToken tries to create a human readable slug.
func (g *HumanReadable) NewToken() (string, error) {
	return g.client.HumanReadableSlug()
}
