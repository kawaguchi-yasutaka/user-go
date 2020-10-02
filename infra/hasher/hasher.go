package hasher

import (
	"golang.org/x/crypto/bcrypt"
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type Hasher struct {
}

func NewHahser() interfaces.IHasher {
	return Hasher{}
}

func (hahser Hasher) GeneratePasswordDigest(password model.UserRawPassword) (model.UserPasswordDigest, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return model.UserPasswordDigest(hashed), nil
}
