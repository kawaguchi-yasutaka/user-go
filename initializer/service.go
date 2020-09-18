package initializer

import "user-go/domain/service"

type Service struct {
	userService service.UserService
}

func NewService(infra Infra, repository Repository) Service {
	return Service{
		userService: service.NewUserService(repository.userRepository, infra.hasher),
	}
}
