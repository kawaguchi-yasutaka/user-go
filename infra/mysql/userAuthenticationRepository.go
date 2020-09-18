package mysql

import "user-go/domain/model"

type UserAuthentication struct {
	UserID         int64  `gorm:"primaryKey"`
	PasswordDigest string `gorm:"not null;default: ''"`
}

func FromUserAuthenticationModel(userPassword model.UserAuthentication) UserAuthentication {
	return UserAuthentication{
		UserID:         int64(userPassword.UserID),
		PasswordDigest: string(userPassword.PasswordDigest),
	}
}
