package mailer

import (
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type mailerMock struct {
}

var _ interfaces.IUserMailer = mailerMock{}

func (m mailerMock) SendActivateCode(to model.UserEmail, code model.UserActivationCode, id model.UserID) error {
	return nil
}

func (m mailerMock) SendMultiAuthenticationCode(to model.UserEmail, code model.UserMultiAuthenticationCode) error {
	panic("not implement")
}
