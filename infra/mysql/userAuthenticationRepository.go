package mysql

import (
	"errors"
	"gorm.io/gorm"
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type UserAuthentication struct {
	UserID                 int64  `gorm:"primaryKey"`
	PasswordDigest         string `gorm:"not null;default: ''"`
	ActivationCode         string `gorm:"not null:default: ''"`
	ActivationCodeExpireAt int64  `gorm: "not null;default: ''`
}

type UserAuthenticationRepository struct {
	db *gorm.DB
}

func NewUserAuthenticationRepository(db *gorm.DB) interfaces.IUserAuthenticationRepository {
	return UserAuthenticationRepository{
		db: db,
	}
}

func FromUserAuthenticationModel(userPassword model.UserAuthentication) UserAuthentication {
	return UserAuthentication{
		UserID:         int64(userPassword.UserID),
		PasswordDigest: string(userPassword.PasswordDigest),
	}
}

func Tomodel(authentication UserAuthentication) model.UserAuthentication {
	return model.UserAuthentication{
		UserID:                 model.UserID(authentication.UserID),
		PasswordDigest:         model.UserPasswordDigest(authentication.PasswordDigest),
		ActivationCode:         model.UserActivationCode(authentication.ActivationCode),
		ActivationCodeExpireAt: model.UserActivationCodeExpireAt(authentication.ActivationCodeExpireAt),
	}
}

func (repo UserAuthenticationRepository) Save(authentication model.UserAuthentication) error {
	a := FromUserAuthenticationModel(authentication)
	if err := repo.db.Save(&a).Error; err != nil {
		return err
	}
	return nil
}

func (repo UserAuthenticationRepository) FindByUserID(UserID model.UserID) (model.UserAuthentication, error) {
	auth := model.UserAuthentication{}
	if result := repo.db.First(&auth, UserID); !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.UserAuthentication{}, result.Error
	}
	return auth, nil
}
