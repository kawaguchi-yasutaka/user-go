package interfaces

import "user-go/domain/model"

type IMailer interface {
	SendAuthenticationCode(to model.UserEmail, code model.UserActivationCode) error
}
