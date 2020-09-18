package initializer

import (
	"user-go/domain/interfaces"
	"user-go/infra/hasher"
)

type Infra struct {
	hasher interfaces.IHasher
}

func NewInfra() Infra {
	return Infra{
		hasher: hasher.NewHahser(),
	}
}
