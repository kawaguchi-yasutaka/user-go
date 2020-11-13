package mysql

import (
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type UserRememberRepositoryMock struct {
	UserRemembers map[model.UserSessionId]model.UserRemember
}

func NewUserRememberRepositoryMock() UserRememberRepositoryMock {
	return UserRememberRepositoryMock{
		UserRemembers: map[model.UserSessionId]model.UserRemember{},
	}
}

var _ interfaces.IUserRemenberRepository = UserRememberRepositoryMock{}

func (r UserRememberRepositoryMock) Save(userRemember model.UserRemember) error {
	r.UserRemembers[userRemember.SessionId] = userRemember
	return nil
}

func (r UserRememberRepositoryMock) FindBySessionId(sessionId model.UserSessionId) (model.UserRemember, error) {
	panic("not implement")
}
