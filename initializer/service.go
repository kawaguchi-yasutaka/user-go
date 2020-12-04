package initializer

import (
	"user-go/domain/service"
	"user-go/domainClient/jwtGeneratorClient"
	"user-go/domainClient/jwtHandlerClient"
)

type Service struct {
	UserService service.UserService
}

func NewService(infra Infra, repository Repository) Service {

	jwtGenerator := jwtGeneratorClient.NewJwtGeneratorClient(infra.jwtGenerator)
	jwtHandler := jwtHandlerClient.NewJwtHandlerClient(infra.jwtHandler, infra.timeKeeper)

	return Service{
		UserService: service.NewUserService(
			repository.userRepository,
			repository.userAuthenticationRepository,
			repository.userRememberRepository,
			infra.hasher,
			infra.userMailer,
			infra.randGenerator,
			infra.timeKeeper,
			jwtGenerator,
			jwtHandler,
		),
	}
}
