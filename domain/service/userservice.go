package service

import (
	"fmt"
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

func (service UserService) Create(email model.UserEmail, password model.UserRawPassword) (model.UserID, error) {
	pDigest, err := service.hasher.GeneratePasswordDigest(password)
	if err != nil {
		return model.UserID(0), err
	}
	userId, err := service.userRepository.Create(model.NewUser(email), pDigest)
	if err != nil {
		return model.UserID(0), err
	}
	code, expiresAt, err := model.NewAuthenticationCode()
	if err != nil {
		return model.UserID(0), err
	}
	auth, err := service.userAuthenticationRepository.FindByUserID(userId)
	if err != nil {
		return model.UserID(0), err
	}

	auth.UpdateActivationCode(code, expiresAt)
	if err := service.userAuthenticationRepository.Save(auth); err != nil {
		return model.UserID(0), err
	}
	return userId, service.userMailer.SendAuthenticationCode(email, code)
}

func (service UserService) Activate(code model.UserActivationCode) error {
	auth, err := service.userAuthenticationRepository.FindByActivateCode(code)
	if err != nil {
		return err
	}
	if model.IsActivationCodeExpired(auth) {
		return fmt.Errorf("code expired")
	}
	user, err := service.userRepository.FindById(auth.UserID)
	if err != nil {
		return err
	}
	if user.IsActivated() {
		return fmt.Errorf("already activated")
	}

	user.Activate()

	return service.userRepository.Save(user)
}

func (service UserService) ReSendActivateCodeEmail(userID model.UserID) error {
	return nil
}
