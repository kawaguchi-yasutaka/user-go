package mysql

import (
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type UserAuthenticationRepositoryMock struct {
	UserAuthentications map[model.UserID]model.UserAuthentication
}

var _ interfaces.IUserAuthenticationRepository = UserAuthenticationRepositoryMock{}

func (r UserAuthenticationRepositoryMock) Save(authentication model.UserAuthentication) error {
	r.UserAuthentications[authentication.UserID] = authentication
	return nil
}

func (r UserAuthenticationRepositoryMock) FindByUserID(UserID model.UserID) (model.UserAuthentication, error) {
	if a, ok := r.UserAuthentications[UserID]; ok {
		return a, nil
	}
	return model.UserAuthentication{}, model.UserAuthenticationNotFound()
}

func (r UserAuthenticationRepositoryMock) FindByActivateCodeAndUserID(code model.UserActivationCode, id model.UserID) (model.UserAuthentication, error) {
	for _, v := range r.UserAuthentications {
		if v.UserID == id && v.ActivationCode == code {
			return v, nil
		}
	}
	return model.UserAuthentication{}, model.UserAuthenticationNotFound()
}
