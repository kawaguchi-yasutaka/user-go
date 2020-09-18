package interfaces

import "user-go/domain/model"

type IHasher interface {
	GeneratePasswordDigest(password model.UserRawPassword) (model.UserPasswordDigest, error)
}
