package hasher

import (
	"fmt"
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type HasherMock struct {
}

var _ interfaces.IHasher = HasherMock{}

func (h HasherMock) GeneratePasswordDigest(password model.UserRawPassword) (model.UserPasswordDigest, error) {
	return model.UserPasswordDigest(password), nil
}

func (h HasherMock) ValidatePassword(password model.UserRawPassword, pDigest model.UserPasswordDigest) error {
	if string(password) != string(pDigest) {
		return model.IncorrectUserPassword(fmt.Sprintf("%v is incorect ", password))
	}
	return nil
}
