package service

import (
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type UserService struct {
	userRepository interfaces.IUserRepository
	hasher         interfaces.IHasher
}

func NewUserService(userRepository interfaces.IUserRepository, hasher interfaces.IHasher) UserService {
	return UserService{
		userRepository: userRepository,
		hasher:         hasher,
	}
}

func (service UserService) Create(email model.UserEmail, password model.UserRawPassword) error {
	pDigest, err := service.hasher.GeneratePasswordDigest(password)
	if err != nil {
		return err
	}
	return service.userRepository.Create(model.NewUser(email), pDigest)
}
