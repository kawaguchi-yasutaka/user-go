package userMailer

import (
	"fmt"
	"user-go/domain/model"
	"user-go/infra/mailer"
)

type UserMailer struct {
	url    string
	mailer mailer.Mailer
}

func NewUserMailer(url string, mailer mailer.Mailer) UserMailer {
	return UserMailer{
		url:    url,
		mailer: mailer,
	}
}

func (m UserMailer) SendActivateCode(
	to model.UserEmail,
	code model.UserActivationCode,
	id model.UserID,
) error {
	body := fmt.Sprintf("認証コードです \n %v/users/%v/activate?code=%v", m.url, id, code)
	return m.mailer.Send([]string{string(to)}, []byte(body))
}

func (m UserMailer) SendMultiAuthenticationCode(
	to model.UserEmail,
	code model.UserMultiAuthenticationCode,
	id model.UserID,
) error {
	body := fmt.Sprintf("2段階認証コードです \n %v/users/%v/multi-authenticate?code=%v", m.url, id, code)
	return m.mailer.Send([]string{string(to)}, []byte(body))
}
