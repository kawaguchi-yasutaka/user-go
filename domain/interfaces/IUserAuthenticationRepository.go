package interfaces

import "user-go/domain/model"

type IUserAuthenticationRepository interface {
	Save(authentication model.UserAuthentication) error
	FindByUserID(userID model.UserID) (model.UserAuthentication, error)
	FindByActivateCode(code model.UserActivationCode) (model.UserAuthentication, error)
	FindBySessionId(sessionId model.UserSessionId) (model.UserAuthentication, error)
	FindByMultiAuthenticateCode(code model.UserMultiAuthenticationCode) (model.UserAuthentication, error)
}

//ユーザー有効化に利用するコードと、二段階認証に利用するコードどうするか。
