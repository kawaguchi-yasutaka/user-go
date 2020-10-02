package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	auth smtp.Auth
	host string
	from string
}

func NewMailer(email string, password string, host string) Mailer {
	return Mailer{
		auth: smtp.PlainAuth("", email, password, host),
		host: host,
		from: email,
	}
}

func (m Mailer) Send(to []string, body []byte) error {
	return smtp.SendMail(fmt.Sprintf("%s:587", m.host), m.auth, m.from, to, body)
}
