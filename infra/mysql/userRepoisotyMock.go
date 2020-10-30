package mysql

import (
	"errors"
	"user-go/domain/interfaces"
	"user-go/domain/model"
	"user-go/lib/myerror"
)

type userRepositoryMock struct {
	users               map[model.UserID]model.User
	userAuthentications map[model.UserID]model.UserAuthentication
}

var _ interfaces.IUserRepository = userRepositoryMock{}

func (r userRepositoryMock) Create(
	user model.User,
	userPassword model.UserPasswordDigest,
) (model.UserID, error) {
	if _, ok := r.users[user.ID]; ok {
		return 0, myerror.DBError(errors.New(myerror.ErrorDBDuplicateEntry))
	}
	if _, ok := r.userAuthentications[user.ID]; ok {
		return 0, myerror.DBError(errors.New(myerror.ErrorDBDuplicateEntry))
	}
	r.users[user.ID] = user
	r.userAuthentications[user.ID] = model.UserAuthentication{UserID: user.ID, PasswordDigest: userPassword}
	return user.ID, nil
}

func (r userRepositoryMock) Save(user model.User) error {
	panic("not implement")
}

func (r userRepositoryMock) FindById(id model.UserID) (model.User, error) {
	panic("not implement")
}

func (r userRepositoryMock) FindByEmail(email model.UserEmail) (model.User, error) {
	panic("not implement")
}
