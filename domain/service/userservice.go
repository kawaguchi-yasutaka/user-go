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
	if err := auth.ValidateActivationCodeExpired(); err != nil {
		return err
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

func (service UserService) Login(email model.UserEmail, password model.UserRawPassword) (model.UserSessionId, error) {
	user, err := service.userRepository.FindByEmail(email)
	if err != nil {
		return model.UserSessionId(""), err
	}
	auth, err := service.userAuthenticationRepository.FindByUserID(user.ID)
	if err != nil {
		return model.UserSessionId(""), err
	}
	if err := service.hasher.ValidatePassword(password, auth.PasswordDigest); err != nil {
		return model.UserSessionId(""), err
	}
	id, expiresAt, err := model.NewUserSessionId()
	if err != nil {
		return model.UserSessionId(0), err
	}
	auth.UpdateSessionInfo(id, expiresAt)
	return id, service.userAuthenticationRepository.Save(auth)
}

func (service UserService) Logind(sessionId model.UserSessionId) (model.UserID, error) {
	auth, err := service.userAuthenticationRepository.FindBySessionId(sessionId)
	if err != nil {
		return 0, model.UserUnauthorized(err.Error())
	}
	if err := auth.ValidateSessionIdExpired(); err != nil {
		return 0, err
	}
	return auth.UserID, nil
}

func (service UserService) ReSendActivateCodeEmail(userID model.UserID) error {
	return nil
}
