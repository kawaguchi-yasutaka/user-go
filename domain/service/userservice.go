package service

import (
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type userService struct {
	userRepository interfaces.IUserRepository
}

func NewUserService(userRepository interfaces.IUserRepository) userService {
	return userService{
		userRepository: userRepository,
	}
}

func (service userService) Create(email model.UserEmail, password model.UserPassword) error {
	pDigest, err := model.NewPasswordDigest(password)
	if err != nil {
		return err
	}
	return service.userRepository.Create(model.NewUser(email, pDigest))
}
