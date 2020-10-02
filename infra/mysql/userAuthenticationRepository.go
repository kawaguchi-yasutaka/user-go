package mysql

import (
	"gorm.io/gorm"
	"user-go/domain/interfaces"
	"user-go/domain/model"
)

type UserAuthentication struct {
	UserID                 int64  `gorm:"primaryKey"`
	PasswordDigest         string `gorm:"not null;default:''"`
	ActivationCode         string `gorm:"not null;default:''"`
	ActivationCodeExpireAt int64  `gorm:"not null;default:0"`
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
		UserID:                 int64(userPassword.UserID),
		PasswordDigest:         string(userPassword.PasswordDigest),
		ActivationCode:         string(userPassword.ActivationCode),
		ActivationCodeExpireAt: int64(userPassword.ActivationCodeExpiresAt),
	}
}

func (authentication UserAuthentication) ToModel() model.UserAuthentication {
	return model.UserAuthentication{
		UserID:                  model.UserID(authentication.UserID),
		PasswordDigest:          model.UserPasswordDigest(authentication.PasswordDigest),
		ActivationCode:          model.UserActivationCode(authentication.ActivationCode),
		ActivationCodeExpiresAt: model.UserActivationCodeExpiresAt(authentication.ActivationCodeExpireAt),
	}
}

func (repo UserAuthenticationRepository) Save(authentication model.UserAuthentication) error {
	a := FromUserAuthenticationModel(authentication)
	if err := repo.db.Save(&a).Debug().Error; err != nil {
		return err
	}
	return nil
}

func (repo UserAuthenticationRepository) FindByUserID(
	userID model.UserID,
) (model.UserAuthentication, error) {
	auth := UserAuthentication{}
	//見つからない場合とdbのエラーを区別していない
	if result := repo.db.Where("user_id = ?", int64(userID)).First(&auth); result.Error != nil {
		return model.UserAuthentication{}, result.Error
	}
	return auth.ToModel(), nil
}

func (repo UserAuthenticationRepository) FindByActivateCode(
	code model.UserActivationCode,
) (model.UserAuthentication, error) {
	auth := UserAuthentication{}
	if result := repo.db.Where("activation_code = ?", string(code)).First(&auth); result.Error != nil {
		return model.UserAuthentication{}, result.Error
	}
	return auth.ToModel(), nil
}
