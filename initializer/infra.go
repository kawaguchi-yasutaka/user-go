package initializer

import (
	"user-go/config"
	"user-go/domain/interfaces"
	"user-go/domainClient/randGenerator"
	"user-go/domainClient/userMailer"
	"user-go/infra/hasher"
	"user-go/infra/jwtgenerator"
	"user-go/infra/jwthandler"
	"user-go/infra/mailer"
	"user-go/infra/timekeeper"
)

type Infra struct {
	hasher        interfaces.IHasher
	userMailer    interfaces.IUserMailer
	randGenerator interfaces.IRandGenerator
	timeKeeper    interfaces.ITimeKeeper
	jwtGenerator  interfaces.IJwtGenerator
	jwtHandler    interfaces.IJwtHandler
}

func NewInfra(config config.Config) Infra {
	m := mailer.NewMailer(config.APP.EmailAddress, config.APP.EmailPassword, config.APP.EmailHost)
	return Infra{
		hasher:        hasher.NewHahser(),
		userMailer:    userMailer.NewUserMailer(config.APP.URL, m),
		randGenerator: randGenerator.NewRandGenerator(),
		timeKeeper:    timekeeper.NewTimeKeeper(),
		jwtGenerator:  jwtgenerator.NewJwtGenerator(config.APP.JwtPrivateKey),
		jwtHandler:    jwthandler.NewJwtHandler(config.APP.JwtPublicKey),
	}
}
