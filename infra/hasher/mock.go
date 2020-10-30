package hasher

import (
	"fmt"
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type hasherMock struct {
}

var _ interfaces.IHasher = hasherMock{}

func (h hasherMock) GeneratePasswordDigest(password model.UserRawPassword) (model.UserPasswordDigest, error) {
	return model.UserPasswordDigest(password), nil
}

func (h hasherMock) ValidatePassword(password model.UserRawPassword, pDigest model.UserPasswordDigest) error {
	if string(password) != string(pDigest) {
		return model.IncorrectUserPassword(fmt.Sprintf("%v is incorect ", password))
	}
	return nil
}
