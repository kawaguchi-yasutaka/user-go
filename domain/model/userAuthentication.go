package model

import (
	"encoding/base64"
	"math/rand"
	"time"
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

func NewUserAuthentication(userId UserID, passwordDigest UserPasswordDigest) UserAuthentication {
	return UserAuthentication{
		UserID:         userId,
		PasswordDigest: passwordDigest,
	}
}

func NewAuthenticationCode() (UserActivationCode, UserActivationCodeExpiresAt, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return UserActivationCode(""), UserActivationCodeExpiresAt(0), err
	}
	return UserActivationCode(
			base64.URLEncoding.EncodeToString(b)),
		UserActivationCodeExpiresAt(unixtime.NewUnixTime(time.Now().Add(time.Duration(24) * time.Hour))),
		nil
}

func (authentication *UserAuthentication) UpdateActivationCode(code UserActivationCode, expiresAt UserActivationCodeExpiresAt) {
	authentication.ActivationCode = code
	authentication.ActivationCodeExpiresAt = expiresAt
}

func NewUserRawPassword(password string) (UserRawPassword, error) {
	if err := Validate.Var(password, "required,min=8"); err != nil {
		return "", err
	}
	return UserRawPassword(password), nil
}
