package service

import (
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type UserService struct {
	userRepository               interfaces.IUserRepository
	userAuthenticationRepository interfaces.IUserAuthenticationRepository
	hasher                       interfaces.IHasher
	mailer                       interfaces.IMailer
}

func NewUserService(
	userRepository interfaces.IUserRepository,
	userAuthenticationRepository interfaces.IUserAuthenticationRepository,
	hasher interfaces.IHasher,
) UserService {
	return UserService{
		userRepository:               userRepository,
		userAuthenticationRepository: userAuthenticationRepository,
		hasher:                       hasher,
	}
}

func (service UserService) Create(email model.UserEmail, password model.UserRawPassword) error {
	pDigest, err := service.hasher.GeneratePasswordDigest(password)
	if err != nil {
		return err
	}
	if err := service.userRepository.Create(model.NewUser(email), pDigest); err != nil {
		return err
	}
	code, expire_at, err := model.NewAuthenticationCode()
	if err != nil {
		return err
	}
	auth, err := service.userAuthenticationRepository.FindByUserID(0)
	if err != nil {
		return err
	}

	auth.UpdateActivationCode(code, expire_at)
	if err := service.userAuthenticationRepository.Save(auth); err != nil {
		return err
	}
	return service.mailer.SendAuthenticationCode(email, code)
}

func (service UserService) Activate() error {
	return nil
}

func (service UserService) ReSendActivateCodeEmail(userID model.UserID) error {
	return nil
}
