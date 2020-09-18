package hasher

import (
	"golang.org/x/crypto/bcrypt"
	"user-go/domain/model"
)

type Hasher struct {
}

func GeneratePasswordDigest(password model.UserRawPassword) (model.UserPasswordDigest, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return model.UserPasswordDigest(hashed), nil
}
