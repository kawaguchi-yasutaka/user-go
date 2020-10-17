package service

import (
	"fmt"
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type UserService struct {
	userRepository               interfaces.IUserRepository
	userAuthenticationRepository interfaces.IUserAuthenticationRepository
	userRememberRepository       interfaces.IUserRemenberRepository
	hasher                       interfaces.IHasher
	userMailer                   interfaces.IUserMailer
}

func NewUserService(
	userRepository interfaces.IUserRepository,
	userAuthenticationRepository interfaces.IUserAuthenticationRepository,
	userRememberRepository interfaces.IUserRemenberRepository,
	hasher interfaces.IHasher,
	userMailer interfaces.IUserMailer,
) UserService {
	return UserService{
		userRepository:               userRepository,
		userAuthenticationRepository: userAuthenticationRepository,
		userRememberRepository:       userRememberRepository,
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
	code, codeExpiresAt, err := model.NewMultiAuthenticationCode()
	if err != nil {
		return model.UserSessionId(""), err
	}
	sessionId, sessionIdExpiresAt, err := model.NewUserSessionId()
	if err != nil {
		return model.UserSessionId(""), err
	}

	userRemember := model.NewUserRememberBySingleFactorAuthentication(
		user.ID,
		sessionId,
		sessionIdExpiresAt,
		code,
		codeExpiresAt,
	)
	if err := service.userRememberRepository.Save(userRemember); err != nil {
		return model.UserSessionId(""), err
	}
	if err := service.userMailer.SendMultiAuthenticationCode(email, code); err != nil {
		return model.UserSessionId(""), err
	}
	return sessionId, nil
}

func (service UserService) Logind(sessionId model.UserSessionId) (model.UserID, error) {
	userRemember, err := service.userRememberRepository.FindBySessionId(sessionId)
	if err != nil {
		return 0, model.UserUnauthorized(err.Error())
	}
	if err := userRemember.ValidateSession(); err != nil {
		return 0, err
	}
	return userRemember.UserID, nil
}

func (service UserService) MultiAuthenticate(
	code model.UserMultiAuthenticationCode,
	sessionId model.UserSessionId,
) error {
	userRemember, err := service.userRememberRepository.FindBySessionId(sessionId)
	if err != nil {
		return err
	}
	if err := userRemember.ValidateMultiAuthenticationCode(code); err != nil {
		return err
	}
	userRemember.Completed()

	return service.userRememberRepository.Save(userRemember)
}

func (service UserService) ReSendActivateCodeEmail(userID model.UserID) error {
	return nil
}
