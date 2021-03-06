package interfaces

import "user-go/domain/model"

type IUserAuthenticationRepository interface {
	Save(authentication model.UserAuthentication) error
	FindByUserID(userID model.UserID) (model.UserAuthentication, error)
	FindByActivateCodeAndUserID(code model.UserActivationCode, id model.UserID) (model.UserAuthentication, error)
}
