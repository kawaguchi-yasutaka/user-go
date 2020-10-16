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
	return userId, service.userMailer.SendActivateCode(email, code, userId)
}

func (service UserService) Activate(code model.UserActivationCode, id model.UserID) error {
	auth, err := service.userAuthenticationRepository.FindByActivateCode(code, id)
	if err != nil {
		return err
	}
	if err := auth.ValidateActivationCodeExpired(); err != nil {
		return err
	}
	user, err := service.userRepository.FindById(id)
	if err != nil {
		return err
	}
	if user.IsActivated() {
		return fmt.Errorf("already activated")
	}

	user.Activate()

	return service.userRepository.Save(user)
}

func (service UserService) Login(email model.UserEmail, password model.UserRawPassword) error {
	user, err := service.userRepository.FindByEmail(email)
	if err != nil {
		return err
	}
	auth, err := service.userAuthenticationRepository.FindByUserID(user.ID)
	if err != nil {
		return err
	}
	if err := service.hasher.ValidatePassword(password, auth.PasswordDigest); err != nil {
		return err
	}
	code, expiresAt, err := model.NewMultiAuthenticationCode()
	if err != nil {
		return err
	}

	auth.UpdateMultiAuthenticationInfo(code, expiresAt)

	if err := service.userAuthenticationRepository.Save(auth); err != nil {
		return err
	}
	return service.userMailer.SendMultiAuthenticationCode(email, code, user.ID)
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

func (service UserService) MultiAuthenticate(
	code model.UserMultiAuthenticationCode,
	id model.UserID,
) (model.UserSessionId, error) {
	auth, err := service.userAuthenticationRepository.FindByMultiAuthenticateCode(code, id)
	if err != nil {
		return model.UserSessionId(""), err
	}
	if err := auth.ValidateMultiAuthenticationCodeExpired(); err != nil {
		return model.UserSessionId(""), err
	}
	sessionId, expiresAt, err := model.NewUserSessionId()
	if err != nil {
		return model.UserSessionId(0), err
	}
	auth.UpdateSessionInfo(sessionId, expiresAt)
	return sessionId, service.userAuthenticationRepository.Save(auth)
}

func (service UserService) ReSendActivateCodeEmail(userID model.UserID) error {
	return nil
}
