package generator

import "github.com/google/uuid"

type UUIDGenerator struct{}

func NewUUIDGenerator() *UUIDGenerator {
	return new(UUIDGenerator)
}

func (g *UUIDGenerator) NewToken() string {
	return uuid.NewString()
}
