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

func (m UserMailer) SendAuthenticationCode(to model.UserEmail, code model.UserActivationCode) error {
	body := fmt.Sprintf("認証コードです \n %v/activate?code=%v", m.url, code)
	return m.mailer.Send([]string{string(to)}, []byte(body))
}
