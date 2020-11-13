package interfaces

import "user-go/domain/model"

type IUserRepository interface {
	Create(user model.User, userPassword model.UserPasswordDigest) (model.UserID, error)
	Save(user model.User) error
	FindById(id model.UserID) (model.User, error)
	FindByEmail(email model.UserEmail) (model.User, error)
}
