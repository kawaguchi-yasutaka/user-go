package interfaces

type IRandGenerator interface {
	GenerateRandByte(size int) ([]byte, error)
}
