package interfaces

import "user-go/domain/model"

type IUserRepository interface {
	Create(user model.User, userPassword model.UserPasswordDigest) error
}
