package mysql

import "user-go/domain/model"

type UserPassword struct {
	UserID         int64  `gorm:"primaryKey"`
	PasswordDigest string `gorm:"not null;default: ''"`
}

func FromUserPasswordModel(userPassword model.UserPassword) UserPassword {
	return UserPassword{
		UserID:         int64(userPassword.UserID),
		PasswordDigest: string(userPassword.PasswordDigest),
	}
}
