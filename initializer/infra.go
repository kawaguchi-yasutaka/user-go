package initializer

import (
	"user-go/config"
	"user-go/domain/interfaces"
	"user-go/domainClient/userMailer"
	"user-go/infra/hasher"
	"user-go/infra/mailer"
)

type Infra struct {
	hasher     interfaces.IHasher
	userMailer interfaces.IUserMailer
}

func NewInfra(config config.Config) Infra {
	m := mailer.NewMailer(config.APP.EmailAddress, config.APP.EmailPassword, config.APP.EmailHost)
	return Infra{
		hasher:     hasher.NewHahser(),
		userMailer: userMailer.NewUserMailer(config.APP.URL, m),
	}
}
