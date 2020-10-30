package mysql

import (
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type userRememberRepositoryMock struct {
	userRemembers map[model.UserID]model.UserRemember
}

var _ interfaces.IUserRemenberRepository = userRememberRepositoryMock{}

func (r userRememberRepositoryMock) Save(userRemember model.UserRemember) error {
	panic("not implement")
}

func (r userRememberRepositoryMock) FindBySessionId(sessionId model.UserSessionId) (model.UserRemember, error) {
	panic("not implement")
}
