package model

import (
	"encoding/base64"
	"math/rand"
	"time"
	"user-go/lib/unixtime"
)

type UserAuthentication struct {
	UserID                 UserID
	PasswordDigest         UserPasswordDigest
	ActivationCode         UserActivationCode
	ActivationCodeExpireAt UserActivationCodeExpireAt
}

type (
	UserRawPassword            string
	UserPasswordDigest         string
	UserActivationCode         string
	UserActivationCodeExpireAt unixtime.UnixTime
)

func NewUserAuthentication(userId UserID, passwordDigest UserPasswordDigest) UserAuthentication {
	return UserAuthentication{
		UserID:         userId,
		PasswordDigest: passwordDigest,
	}
}

func NewAuthenticationCode() (UserActivationCode, UserActivationCodeExpireAt, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return UserActivationCode(""), UserActivationCodeExpireAt(0), err
	}
	return UserActivationCode(
			base64.URLEncoding.EncodeToString(b)),
		UserActivationCodeExpireAt(unixtime.NewUnixTime(time.Now().Add(time.Duration(24) * time.Hour))),
		nil
}

func (authentication *UserAuthentication) UpdateActivationCode(code UserActivationCode, expireAt UserActivationCodeExpireAt) {
	authentication.ActivationCode = code
	authentication.ActivationCodeExpireAt = expireAt
}

func NewUserRawPassword(password string) (UserRawPassword, error) {
	if err := Validate.Var(password, "required,min=8"); err != nil {
		return "", err
	}
	return UserRawPassword(password), nil
}
