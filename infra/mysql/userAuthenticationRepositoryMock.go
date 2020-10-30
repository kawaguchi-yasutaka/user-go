package mysql

import (
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type userAuthenticationRepositoryMock struct {
	userAuthentications map[model.UserID]model.UserAuthentication
}

var _ interfaces.IUserAuthenticationRepository = userAuthenticationRepositoryMock{}

func (r userAuthenticationRepositoryMock) Save(authentication model.UserAuthentication) error {
	r.userAuthentications[authentication.UserID] = authentication
	return nil
}

func (r userAuthenticationRepositoryMock) FindByUserID(UserID model.UserID) (model.UserAuthentication, error) {
	panic("not implement")
}

func (r userAuthenticationRepositoryMock) FindByActivateCode(code model.UserActivationCode, id model.UserID) (model.UserAuthentication, error) {
	panic("not implement")
}
