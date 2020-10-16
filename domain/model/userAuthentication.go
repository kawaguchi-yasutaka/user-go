package model

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/go-playground/validator/v10"
	mathRand "math/rand"
	"net/http"
	"time"
	"user-go/lib/myerror"
	"user-go/lib/unixtime"
)

type UserAuthentication struct {
	UserID                           UserID
	PasswordDigest                   UserPasswordDigest
	ActivationCode                   UserActivationCode
	ActivationCodeExpiresAt          UserActivationCodeExpiresAt
	MultiAuthenticationCode          UserMultiAuthenticationCode
	MultiAuthenticationCodeExpiresAt UserMultiAuthenticationCodeExpiresAt
	SessionId                        *UserSessionId
	SessionIdExpiresAt               UserSessionIdExpiresAt
}

type (
	UserRawPassword                      string
	UserPasswordDigest                   string
	UserActivationCode                   string
	UserActivationCodeExpiresAt          unixtime.UnixTime
	UserMultiAuthenticationCode          string
	UserMultiAuthenticationCodeExpiresAt unixtime.UnixTime
	UserSessionId                        string
	UserSessionIdExpiresAt               unixtime.UnixTime
)

const (
	ErrorUserUnauthorized                   myerror.ErrorType = "user_unauthorized"
	ErrorUserAuthenticationNotFound         myerror.ErrorType = "User_authentication_not_found"
	ErrorRequiredUserPassword               myerror.ErrorType = "required_user_password"
	ErrorTooShortUserPassword               myerror.ErrorType = "too_short_user_password"
	ErrorExpiredUserActivationCode          myerror.ErrorType = "expired_user_activation_code"
	ErrorExpiredUserMultiAuthenticationCode myerror.ErrorType = "expired_user_multi_authentication_code"
)

func ExpiredUserActivationCode() myerror.CustomError {
	return myerror.NewCustomError("expired activation code", ErrorExpiredUserActivationCode, http.StatusBadRequest)
}

func UserUnauthorized(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorUserUnauthorized, http.StatusUnauthorized)
}

func UserAuthenticationNotFound() myerror.CustomError {
	return myerror.NewCustomError("not found", ErrorUserAuthenticationNotFound, http.StatusNotFound)
}

func RequiredUserPassword(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorRequiredUserPassword, http.StatusBadRequest)
}

func TooShortUserPassword(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorTooShortUserPassword, http.StatusBadRequest)
}

func ExpiredUserMultiAuthenticationCode() myerror.CustomError {
	return myerror.NewCustomError("expired multi authentication code", ErrorExpiredUserMultiAuthenticationCode, http.StatusBadRequest)
}

func NewUserAuthentication(userId UserID, passwordDigest UserPasswordDigest) UserAuthentication {
	return UserAuthentication{
		UserID:         userId,
		PasswordDigest: passwordDigest,
	}
}

func NewAuthenticationCode() (UserActivationCode, UserActivationCodeExpiresAt, error) {
	b := make([]byte, 64)
	if _, err := mathRand.Read(b); err != nil {
		return UserActivationCode(""), UserActivationCodeExpiresAt(0), err
	}
	return UserActivationCode(
			base64.URLEncoding.EncodeToString(b)),
		UserActivationCodeExpiresAt(unixtime.NewUnixTime(time.Now().Add(time.Duration(24) * time.Hour))),
		nil
}

func NewUserRawPassword(password string) (UserRawPassword, error) {
	if err := Validate.Var(password, "required,min=8"); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			switch fieldName {
			case "required":
				return "", RequiredUserPassword("")
			case "min=8":
				return "", TooShortUserPassword("")
			}
		}
		return "", err
	}
	return UserRawPassword(password), nil
}

func NewMultiAuthenticationCode() (UserMultiAuthenticationCode, UserMultiAuthenticationCodeExpiresAt, error) {
	b := make([]byte, 64)
	if _, err := mathRand.Read(b); err != nil {
		return UserMultiAuthenticationCode(""), UserMultiAuthenticationCodeExpiresAt(0), err
	}
	return UserMultiAuthenticationCode(
			base64.URLEncoding.EncodeToString(b)),
		UserMultiAuthenticationCodeExpiresAt(unixtime.NewUnixTime(time.Now().Add(time.Duration(24) * time.Hour))),
		nil
}

func NewUserSessionId() (UserSessionId, UserSessionIdExpiresAt, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return UserSessionId(""), UserSessionIdExpiresAt(0), err
	}
	return UserSessionId(
			base64.URLEncoding.EncodeToString(b)),
		UserSessionIdExpiresAt(unixtime.NewUnixTime(time.Now().Add(time.Duration(24) * time.Hour))),
		nil
}

func (authentication *UserAuthentication) UpdateActivationCode(code UserActivationCode, expiresAt UserActivationCodeExpiresAt) {
	authentication.ActivationCode = code
	authentication.ActivationCodeExpiresAt = expiresAt
}

func (authentication *UserAuthentication) UpdateSessionInfo(id UserSessionId, expiresAt UserSessionIdExpiresAt) {
	authentication.SessionId = &id
	authentication.SessionIdExpiresAt = expiresAt
}

func (authentication *UserAuthentication) UpdateMultiAuthenticationInfo(code UserMultiAuthenticationCode, expiresAt UserMultiAuthenticationCodeExpiresAt) {
	authentication.MultiAuthenticationCode = code
	authentication.MultiAuthenticationCodeExpiresAt = expiresAt
}

func (authentication UserAuthentication) ValidateActivationCodeExpired() error {
	if unixtime.UnixTime(authentication.ActivationCodeExpiresAt) <= unixtime.Now() {
		return ExpiredUserActivationCode()
	}
	return nil
}

func (authentication UserAuthentication) ValidateMultiAuthenticationCodeExpired() error {
	if unixtime.UnixTime(authentication.MultiAuthenticationCodeExpiresAt) <= unixtime.Now() {
		return ExpiredUserMultiAuthenticationCode()
	}
	return nil
}

func (authentication UserAuthentication) ValidateSessionIdExpired() error {
	if unixtime.UnixTime(authentication.SessionIdExpiresAt) <= unixtime.Now() {
		return UserUnauthorized("session id is expired")
	}
	return nil
}
