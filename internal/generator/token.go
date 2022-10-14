package generator

import "strings"

var builtinNouns = []string{
	"tornado",
	"mahogany",
	"boldness",
	"fan",
	"majority",
	"panda",
	"dragon",
}

var builtinAdjectives = []string{
	"criminal",
	"adoring",
	"vacant",
	"mysterious",
	"frightful",
	"hungry",
	"sad",
}

var builtinColors = []string{
	"violet",
	"indigo",
	"blue",
	"green",
	"yellow",
	"orange",
	"red",
}

type RandInt func() int

type TokenGenerator struct {
	randInt RandInt
}

func NewTokenGenerator(g RandInt) *TokenGenerator {
	return &TokenGenerator{
		randInt: g,
	}
}

func (g *TokenGenerator) NewToken() string {
	parts := []string{
		builtinColors[g.randInt()%len(builtinColors)],
		builtinAdjectives[g.randInt()%len(builtinAdjectives)],
		builtinNouns[g.randInt()%len(builtinNouns)],
	}

	return strings.Join(parts, "-")
}
