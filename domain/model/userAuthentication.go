package model

import (
	"encoding/base64"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
	"user-go/lib/myerror"
	"user-go/lib/unixtime"
)

type UserAuthentication struct {
	UserID                  UserID
	PasswordDigest          UserPasswordDigest
	ActivationCode          UserActivationCode
	ActivationCodeExpiresAt UserActivationCodeExpiresAt
}

type (
	UserRawPassword             string
	UserPasswordDigest          string
	UserActivationCode          string
	UserActivationCodeExpiresAt unixtime.UnixTime
)

const (
	ErrorUserUnauthorized           myerror.ErrorType = "user_unauthorized"
	ErrorUserAuthenticationNotFound myerror.ErrorType = "user_authentication_not_found"
	ErrorRequiredUserPassword       myerror.ErrorType = "required_user_password"
	ErrorTooShortUserPassword       myerror.ErrorType = "too_short_user_password"
	ErrorExpiredUserActivationCode  myerror.ErrorType = "expired_user_activation_code"
)

func ExpiredUserActivationCode() myerror.CustomError {
	return myerror.NewCustomError("expired activation code", ErrorExpiredUserActivationCode, http.StatusBadRequest)
}

func UserUnauthorized(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorUserUnauthorized, http.StatusUnauthorized)
}

func UserAuthenticationNotFound() myerror.CustomError {
	return myerror.NewCustomError("authorization not found", ErrorUserAuthenticationNotFound, http.StatusNotFound)
}

func RequiredUserPassword(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorRequiredUserPassword, http.StatusBadRequest)
}

func TooShortUserPassword(msg string) myerror.CustomError {
	return myerror.NewCustomError(msg, ErrorTooShortUserPassword, http.StatusBadRequest)
}

func NewUserAuthentication(userId UserID, passwordDigest UserPasswordDigest) UserAuthentication {
	return UserAuthentication{
		UserID:         userId,
		PasswordDigest: passwordDigest,
	}
}

func NewAuthenticationCode(rand []byte, now unixtime.UnixTime) (UserActivationCode, UserActivationCodeExpiresAt, error) {
	return UserActivationCode(base64.URLEncoding.EncodeToString(rand)),
		UserActivationCodeExpiresAt(now + unixtime.UnixTime(time.Hour*24)),
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

func (authentication *UserAuthentication) UpdateActivationCode(code UserActivationCode, expiresAt UserActivationCodeExpiresAt) {
	authentication.ActivationCode = code
	authentication.ActivationCodeExpiresAt = expiresAt
}

func (authentication UserAuthentication) ValidateActivationCodeExpired(now unixtime.UnixTime) error {
	if unixtime.UnixTime(authentication.ActivationCodeExpiresAt) <= now {
		return ExpiredUserActivationCode()
	}
	return nil
}
