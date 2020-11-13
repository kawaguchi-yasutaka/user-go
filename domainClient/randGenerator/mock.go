package randGenerator

import "user-go/domain/interfaces"

type RandGeneratorMock struct {
	RandByte []byte
}

func NewRandGeneratorMock() RandGeneratorMock {
	return RandGeneratorMock{}
}

var _ interfaces.IRandGenerator = RandGeneratorMock{}

func (r RandGeneratorMock) GenerateRandByte(size int) ([]byte, error) {
	return r.RandByte, nil
}
