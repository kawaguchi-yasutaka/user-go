package randGenerator

import (
	"math/rand"
	"user-go/domain/interfaces"
)

type RandGenerator struct {
}

var _ interfaces.IRandGenerator = RandGenerator{}

func NewRandGenerator() RandGenerator {
	return RandGenerator{}
}

func (r RandGenerator) GenerateRandByte(size int) ([]byte, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
