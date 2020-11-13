package userMailer

import (
	"user-go/domain/model"
)

type UserMailerMock struct {
	ActivateCodeSendedCount            *int
	MultiAuthenticationCodeSendedCount *int
}

func NewUserMailerMock() UserMailerMock {
	aCount := 0
	mCount := 0
	return UserMailerMock{
		ActivateCodeSendedCount:            &aCount,
		MultiAuthenticationCodeSendedCount: &mCount,
	}
}

func (m UserMailerMock) SendActivateCode(
	to model.UserEmail,
	code model.UserActivationCode,
	id model.UserID,
) error {
	*m.ActivateCodeSendedCount += 1
	return nil
}

func (m UserMailerMock) SendMultiAuthenticationCode(
	to model.UserEmail,
	code model.UserMultiAuthenticationCode,
) error {
	*m.MultiAuthenticationCodeSendedCount += 1
	return nil
}
