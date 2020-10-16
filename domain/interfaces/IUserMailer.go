package interfaces

import "user-go/domain/model"

type IUserMailer interface {
	SendActivateCode(to model.UserEmail, code model.UserActivationCode) error
	SendMultiAuthenticationCode(to model.UserEmail, code model.UserMultiAuthenticationCode) error
}
