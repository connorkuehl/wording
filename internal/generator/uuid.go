package generator

import (
	"github.com/google/uuid"
)

// UUIDGenerator generates UUIDs.
type UUIDGenerator struct{}

// NewUUIDGenerator creates a new UUIDGenerator.
func NewUUIDGenerator() *UUIDGenerator {
	return new(UUIDGenerator)
}

// NewToken creates a UUID token.
func (g *UUIDGenerator) NewToken() string {
	return uuid.NewString()
}
