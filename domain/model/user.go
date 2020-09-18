package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	Email          UserEmail
	PasswordDigest UserPasswordDigest
	Status         UserStatus
}

type (
	UserStatus         string
	UserEmail          string
	UserPasswordDigest string
)

type UserPassword string

const (
	UserStatusInitialized UserStatus = "initialized"
	UserStatusActivated   UserStatus = "activated" //認証済み
)

func NewUser(email UserEmail, digest UserPasswordDigest) User {
	return User{
		Email:          email,
		PasswordDigest: digest,
		Status:         UserStatusInitialized,
	}
}

func NewPasswordDigest(password UserPassword) (UserPasswordDigest, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return UserPasswordDigest(hashed), nil
}
