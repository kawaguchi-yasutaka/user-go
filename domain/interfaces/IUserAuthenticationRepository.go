package interfaces

import "user-go/domain/model"

type IUserAuthenticationRepository interface {
	Save(authentication model.UserAuthentication) error
	FindByUserID(userID model.UserID) (model.UserAuthentication, error)
	FindByActivateCode(code model.UserActivationCode, id model.UserID) (model.UserAuthentication, error)
	FindBySessionId(sessionId model.UserSessionId) (model.UserAuthentication, error)
	FindByMultiAuthenticateCode(code model.UserMultiAuthenticationCode, id model.UserID) (model.UserAuthentication, error)
}
