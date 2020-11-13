package mysql

import (
	"errors"
	"user-go/domain/interfaces"
	"user-go/domain/model"
	"user-go/lib/myerror"
)

type UserRepositoryMock struct {
	Users               map[model.UserID]model.User
	UserAuthentications map[model.UserID]model.UserAuthentication
}

var _ interfaces.IUserRepository = UserRepositoryMock{}

func NewUserRepositoryMock() UserRepositoryMock {
	return UserRepositoryMock{
		Users:               map[model.UserID]model.User{},
		UserAuthentications: map[model.UserID]model.UserAuthentication{},
	}
}

func (r UserRepositoryMock) Create(
	user model.User,
	userPassword model.UserPasswordDigest,
) (model.UserID, error) {
	uNextId := model.UserID(1)
	for k, _ := range r.Users {
		if uNextId <= k {
			uNextId = k + 1
		}
	}
	if _, ok := r.Users[uNextId]; ok {
		return 0, myerror.DBError(errors.New(myerror.ErrorDBDuplicateEntry))
	}
	if _, ok := r.UserAuthentications[uNextId]; ok {
		return 0, myerror.DBError(errors.New(myerror.ErrorDBDuplicateEntry))
	}
	user.ID = uNextId
	r.Users[user.ID] = user
	r.UserAuthentications[user.ID] = model.UserAuthentication{UserID: uNextId, PasswordDigest: userPassword}
	return user.ID, nil
}

func (r UserRepositoryMock) Save(user model.User) error {
	if _, ok := r.Users[user.ID]; ok {
		r.Users[user.ID] = user
		return nil
	}
	return myerror.DBError(errors.New(myerror.ErrorDBError))
}

func (r UserRepositoryMock) FindById(id model.UserID) (model.User, error) {
	if user, ok := r.Users[id]; ok {
		return user, nil
	}
	return model.User{}, model.UserNotFound()
}

func (r UserRepositoryMock) FindByEmail(email model.UserEmail) (model.User, error) {
	for _, u := range r.Users {
		if u.Email == email {
			return u, nil
		}
	}
	return model.User{}, model.UserNotFound()
}
