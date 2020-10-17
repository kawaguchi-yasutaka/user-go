package initializer

import (
	"user-go/domain/service"
)

type Service struct {
	UserService service.UserService
}

func NewService(infra Infra, repository Repository) Service {
	return Service{
		UserService: service.NewUserService(
			repository.userRepository,
			repository.userAuthenticationRepository,
			repository.userRememberRepository,
			infra.hasher,
			infra.userMailer,
		),
	}
}
