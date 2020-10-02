package service

import (
	"log"
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type UserService struct {
	userRepository               interfaces.IUserRepository
	userAuthenticationRepository interfaces.IUserAuthenticationRepository
	hasher                       interfaces.IHasher
	userMailer                   interfaces.IUserMailer
}

func NewUserService(
	userRepository interfaces.IUserRepository,
	userAuthenticationRepository interfaces.IUserAuthenticationRepository,
	hasher interfaces.IHasher,
	userMailer interfaces.IUserMailer,
) UserService {
	return UserService{
		userRepository:               userRepository,
		userAuthenticationRepository: userAuthenticationRepository,
		hasher:                       hasher,
		userMailer:                   userMailer,
	}
}

func (service UserService) Create(email model.UserEmail, password model.UserRawPassword) error {
	pDigest, err := service.hasher.GeneratePasswordDigest(password)
	if err != nil {
		return err
	}
	userId, err := service.userRepository.Create(model.NewUser(email), pDigest)
	if err != nil {
		return err
	}
	code, expiresAt, err := model.NewAuthenticationCode()
	if err != nil {
		return err
	}
	auth, err := service.userAuthenticationRepository.FindByUserID(userId)
	if err != nil {
		return err
	}
	log.Printf("code: %v", code)
	log.Printf("expiresAt: %v", expiresAt)

	auth.UpdateActivationCode(code, expiresAt)
	if err := service.userAuthenticationRepository.Save(auth); err != nil {
		return err
	}
	return service.userMailer.SendAuthenticationCode(email, code)
}

func (service UserService) Activate() error {
	return nil
}

func (service UserService) ReSendActivateCodeEmail(userID model.UserID) error {
	return nil
}
