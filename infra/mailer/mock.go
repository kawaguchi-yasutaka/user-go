package mailer

import (
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type MailerMock struct {
}

var _ interfaces.IUserMailer = MailerMock{}

func (m MailerMock) SendActivateCode(to model.UserEmail, code model.UserActivationCode, id model.UserID) error {
	return nil
}

func (m MailerMock) SendMultiAuthenticationCode(to model.UserEmail, code model.UserMultiAuthenticationCode) error {
	panic("not implement")
}
