package interfaces

import "user-go/domain/model"

type IUserMailer interface {
	SendAuthenticationCode(to model.UserEmail, code model.UserActivationCode) error
}
