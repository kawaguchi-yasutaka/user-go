package mysql

import (
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type UserRememberRepositoryMock struct {
	userRemembers map[model.UserID]model.UserRemember
}

var _ interfaces.IUserRemenberRepository = UserRememberRepositoryMock{}

func (r UserRememberRepositoryMock) Save(userRemember model.UserRemember) error {
	panic("not implement")
}

func (r UserRememberRepositoryMock) FindBySessionId(sessionId model.UserSessionId) (model.UserRemember, error) {
	panic("not implement")
}
